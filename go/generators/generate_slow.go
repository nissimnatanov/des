package generators

import (
	"context"
	"time"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/solver"
)

type slowStage struct {
	GeneratePerCandidate int // how many candidates to fork from each
	SelectBest           int
	PreserveAtLeastOne   bool // preserve at least one candidate from each base
	FreeCells            int
	MinToRemove          int
	MaxToRemove          int
}

var slowStages = []slowStage{
	// at first we only have one candidate (e.g. solution), we can create many best forks
	{FreeCells: 40, MinToRemove: 10, MaxToRemove: 15, GeneratePerCandidate: 10, SelectBest: 10},
	{FreeCells: 45, MinToRemove: 2, MaxToRemove: 5, GeneratePerCandidate: 10, SelectBest: 20, PreserveAtLeastOne: true},
	{FreeCells: 50, MinToRemove: 1, MaxToRemove: 3, GeneratePerCandidate: 5, SelectBest: 25, PreserveAtLeastOne: true},
	{FreeCells: 55, MinToRemove: 1, MaxToRemove: 2, GeneratePerCandidate: 10, SelectBest: 10, PreserveAtLeastOne: false},
	// this is the last stage, we can throw away the candidates that did not make the cut
	{FreeCells: solver.MaxFreeCellsForValidBoard, MinToRemove: 1, MaxToRemove: 1, GeneratePerCandidate: 20},
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

	// by default, replace the solution every 10 tries, but with the higher levels
	// prefer to keep it longer to benefit more from the cache
	replaceSolutionEvery := 10
	if g.lr.Min >= solver.LevelNightmare {
		replaceSolutionEvery = 500
	} else if g.lr.Min >= solver.LevelDarkEvil {
		replaceSolutionEvery = 100
	}

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
		for si := 0; si < len(slowStages) && ctx.Err() == nil; si++ {
			stage := slowStages[si]
			newFinal, newCandidates := g.generateSlowStage(ctx, candidates, stage)
			if newFinal.Size() == 0 && (len(newCandidates) == 0 || si == len(slowStages)-1) {
				// if we got no finals and we are at the last stage or no new candidates left,
				// report this stage as empty
				stageStats.Report(0, si)
				break
			}
			if g.onNewResult != nil {
				for _, res := range newFinal.Results() {
					g.onNewResult(res)
				}
			}
			finalCandidates.AddAll(newFinal)
			if newFinal.Size() > 0 {
				stageStats.Report(newFinal.Size(), si)
			}
			if hasEnoughFinalCandidates(finalCandidates, g.count) {
				// stop once we have at least one if count was requested
				break generationLoop
			}
			candidates = combineCandidates(newCandidates, stage.SelectBest, stage.PreserveAtLeastOne)
		}

		if tries%replaceSolutionEvery == 0 {
			// if we are stuck on the same solution for too long, try the next one
			// it also helps reducing the per-solution cache footprint
			cacheStats.MergeAndDrain(initState.SolutionState().Cache().Stats())
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
		switch bs.Candidates().Size() {
		case 0:
			// this board can no longer be enhanced
			if bs.Progress() >= internal.InRangeKeepGoing && bs.Progress() <= internal.InRangeStop {
				finalCandidates.Add(bs)
			}
			// throw this candidate away, it is below the desired level
			continue
		case 1:
			// last candidate, forking is useless - just try to remove it
			bs = bs.RemoveOneByOne(ctx, stage.FreeCells)
			if bs != nil {
				if bs.Progress() == internal.InRangeKeepGoing || bs.Progress() == internal.InRangeStop {
					finalCandidates.Add(bs)
				}
			}
			continue
		}
		for range stage.GeneratePerCandidate {
			if bs.Candidates() == indexes.MinBitSet81 {
				// no more candidates to remove, we can stop
				if bs.Progress() >= internal.InRangeKeepGoing && bs.Progress() <= internal.InRangeStop {
					finalCandidates.Add(bs)
				}
				break
			}
			if bs.Progress() == internal.InRangeStop {
				// we can stop with this candidate if it is already at the desired level
				finalCandidates.Add(bs)
				break
			}
			var bsForked *internal.BoardState
			// can refine more, try random first
			if stage.MaxToRemove > 1 {
				bsForked = bs.Remove(ctx, internal.RemoveArgs{
					FreeCells:        stage.FreeCells,
					BatchMinToRemove: stage.MinToRemove,
					BatchMaxToRemove: stage.MaxToRemove,
					BatchMaxTries:    5,
				})
			}
			if bsForked == nil {
				bsForked = bs.RemoveOneByOne(ctx, stage.FreeCells)
			}
			if bsForked == nil {
				// try another removal combination of the same board
				continue
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
				if newPerBoard.Size() >= stage.SelectBest {
					// we have enough candidates for this board, no need to fork more
					break
				}
			}
		}
	}
	return
}

func combineCandidates(
	newCandidates []*internal.SortedBoardStates,
	selectBest int,
	preserveAtLeastOne bool,
) *internal.SortedBoardStates {
	switch {
	case len(newCandidates) == 0 || selectBest <= 0:
		return internal.NewSortedBoardStates()
	case len(newCandidates) == 1:
		newCandidates[0].TrimSize(selectBest)
		return newCandidates[0]
	}

	sbs := internal.NewSortedBoardStates()
	indexes := make([]int, len(newCandidates))
	if preserveAtLeastOne {
		// preserve at least one best candidate from each base
		for i, bs := range newCandidates {
			if bs.Size() > 0 {
				sbs.Add(bs.Get(0))
				indexes[i]++
			}
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
