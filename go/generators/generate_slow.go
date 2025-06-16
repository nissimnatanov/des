package generators

import (
	"context"
	"slices"
	"time"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/solver"
)

type slowStage struct {
	GeneratePerCandidate int // how many candidates to fork from each
	SelectBest           int
	FreeCells            int
	MinToRemove          int
	MaxToRemove          int
}

var slowStages = []slowStage{
	// at first we only have one candidate (e.g. solution), so we can create many forks
	{FreeCells: 45, MinToRemove: 10, MaxToRemove: 15, GeneratePerCandidate: 20, SelectBest: 10},
	{FreeCells: 55, MinToRemove: 2, MaxToRemove: 4, GeneratePerCandidate: 5, SelectBest: 15},
	{FreeCells: 60, MinToRemove: 1, MaxToRemove: 2, GeneratePerCandidate: 5, SelectBest: 20},
	// this is the last stage, we can throw away the candidates that did not make the cut
	{FreeCells: solver.MaxFreeCellsForValidBoard, MinToRemove: 1, MaxToRemove: 1, GeneratePerCandidate: 5},
}

func hasEnoughFinalCandidates(finalCandidates *internal.SortedBoardStates, requestedCount int) bool {
	if requestedCount <= 0 {
		return finalCandidates.Size() > 0
	}
	return finalCandidates.Size() >= requestedCount
}

func (g *Generator) generateSlow(ctx context.Context) []*solver.Result {
	tries := 0

	candidates := internal.NewSortedBoardStates()
	finalCandidates := internal.NewSortedBoardStates()
	var stageStats internal.GamePerStageStats

	// Start with basic boards in a easy range, so we can generate and sort a lot of
	// candidates. If we start with harder boards, the generation becomes slower.
	start := time.Now()

	// turn on use cache
	initState := g.newInitialBoardState(ctx, true)
	cacheStats := solver.CacheStats{}

generationLoop:
	for ctx.Err() == nil {
		tries++

		candidates.Reset()
		startState := initState
		startState = g.removeSingleValue(ctx, initState)
		candidates.Add(startState)

		// enhance the candidates to the desired level
		for stage := 0; stage < len(slowStages) && ctx.Err() == nil; stage++ {
			newFinal, newCandidates := g.generateSlowStage(ctx, candidates, slowStages[stage])
			if newFinal.Size() == 0 && len(newCandidates) == 0 {
				stageStats.Report(0, stage)
				break
			}
			if g.onNewResult != nil {
				for _, res := range newFinal.Results() {
					g.onNewResult(res)
				}
			}
			finalCandidates.AddAll(newFinal)
			if newFinal.Size() > 0 {
				stageStats.Report(newFinal.Size(), stage)
			}
			if hasEnoughFinalCandidates(finalCandidates, g.count) {
				// stop once we have at least one if count was requested
				break generationLoop
			}
			candidates = combineCandidates(newCandidates, slowStages[stage].SelectBest)
		}

		if tries%25 == 0 {
			// if we are stuck on the same solution for too long, try the next one
			// it also helps reducing the per-solution cache footprint
			// TODO cache: .MergeAndDrain(initState.SolutionState().Cache().Stats())
			initState = g.newInitialBoardState(ctx, true)
		}
	}

	// return the results so far, even if ctx canceled in the middle
	cacheStats.MergeAndDrain(initState.SolutionState().Cache().Stats())
	internal.Stats.ReportGeneration(finalCandidates.Size(), time.Since(start), int64(tries), stageStats, cacheStats)
	if g.count > 0 {
		finalCandidates.TrimSize(g.count)
	}
	return finalCandidates.Results()
}

func (g *Generator) generateSlowStage(
	ctx context.Context,
	candidates *internal.SortedBoardStates,
	stage slowStage,
) (finalCandidates *internal.SortedBoardStates, newCandidates []*internal.SortedBoardStates) {
	finalCandidates = internal.NewSortedBoardStates()
	// refine the candidates to the desired level
	for bs := range candidates.Boards {
		var newPerBoard *internal.SortedBoardStates
		indexCandidates := slices.Collect(bs.Candidates().Indexes)
		switch len(indexCandidates) {
		case 0:
			// this board can no longer be enhanced
			if bs.Progress() >= internal.InRangeKeepGoing && bs.Progress() <= internal.InRangeStop {
				finalCandidates.Add(bs)
			}
			// throw this candidate away, it is below the desired level
			continue
		case 1:
			// last candidate, forking is useless - just try to remove it
			bs = bs.RemoveCandidatesOneByOne(ctx, indexCandidates)
			if bs != nil {
				if bs.Progress() == internal.InRangeKeepGoing || bs.Progress() == internal.InRangeStop {
					finalCandidates.Add(bs)
				}
			}
			continue
		}
		forkCount := stage.GeneratePerCandidate
		for range forkCount {
			if bs.Candidates() == indexes.MinBitSet81 {
				// we can no longer enhance this board, add it to final candidates
				if bs.Progress() >= internal.InRangeKeepGoing && bs.Progress() <= internal.InRangeStop {
					finalCandidates.Add(bs)
				}
				// otherwise stop trying this board
				break
			}
			// can refine more, try random first
			bsForked := bs.Remove(ctx, internal.RemoveArgs{
				FreeCells:        stage.FreeCells,
				BatchMinToRemove: stage.MinToRemove,
				BatchMaxToRemove: stage.MaxToRemove,
				BatchMaxTries:    10,
			})
			if bsForked == nil {
				bsForked = bs.RemoveCandidatesOneByOne(ctx, indexCandidates)
			}
			if bsForked == nil {
				break
			}
			if bsForked.Progress() == internal.InRangeStop ||
				(stage.SelectBest == 0 && bsForked.Progress() == internal.InRangeKeepGoing) {
				// capture if we have reached the desired level, or if we are at level but
				// also the last stage
				finalCandidates.Add(bsForked)
			} else {
				if newPerBoard == nil {
					// first child candidate for the current one
					newPerBoard = internal.NewSortedBoardStates()
					newCandidates = append(newCandidates, newPerBoard)
				}
				newPerBoard.Add(bsForked)
			}
		}
	}
	return
}

func combineCandidates(newCandidates []*internal.SortedBoardStates, selectBest int) *internal.SortedBoardStates {
	switch {
	case len(newCandidates) == 0 || selectBest <= 0:
		return internal.NewSortedBoardStates()
	case len(newCandidates) == 1:
		newCandidates[0].TrimSize(selectBest)
		return newCandidates[0]
	}

	sbs := internal.NewSortedBoardStates()
	indexes := make([]int, len(newCandidates))
	// always preserve at least one best candidate from each base
	for i, bs := range newCandidates {
		if bs.Size() > 0 {
			sbs.Add(bs.Get(0))
			indexes[i]++
		}
	}
	// in case we overflow in the first loop, remove the excess
	sbs.TrimSize(selectBest)
	for sbs.Size() < selectBest {
		var best *internal.BoardState
		var bestIndex int
		for i, bi := range indexes {
			if bi >= newCandidates[i].Size() {
				// no more candidates in this set
				continue
			}
			cur := newCandidates[i].Get(bi)
			if best == nil || best.Complexity() < cur.Complexity() {
				best = cur
				bestIndex = i
			}
		}
		if best == nil {
			break // no more candidates to add
		}
		indexes[bestIndex]++
		sbs.Add(best)
	}
	return sbs
}
