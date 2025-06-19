package generators

import (
	"context"
	"time"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/internal/stats"
	"github.com/nissimnatanov/des/go/solver"
)

type slowStage struct {
	GeneratePerCandidate int // how many candidates to fork from each
	SelectBest           int
	FreeCells            int
	MinToRemove          int
	MaxToRemove          int
	// minimum complexities per desired min level to consider a candidate valid for this stage
	MinComplexities map[solver.Level]solver.StepComplexity
	// TopN means we are interested in the top N sub-candidates for each candidate
	TopN int
}

var slowStages = []slowStage{
	// at first we only have one candidate (e.g. solution), we can create many best forks
	{FreeCells: 42, MinToRemove: 10, MaxToRemove: 15, GeneratePerCandidate: 50, SelectBest: 50},
	{FreeCells: 49, MinToRemove: 2, MaxToRemove: 5, GeneratePerCandidate: 20, SelectBest: 100},
	{FreeCells: 51, TopN: 10, SelectBest: 50},
	{FreeCells: 55, MinToRemove: 2, MaxToRemove: 3, GeneratePerCandidate: 10, SelectBest: 40},
	{FreeCells: 60, MinToRemove: 1, MaxToRemove: 3, GeneratePerCandidate: 10, SelectBest: 30},
	//{FreeCells: solver.MaxFreeCellsForValidBoard, MinToRemove: 1, MaxToRemove: 2, GeneratePerCandidate: 20, SelectBest: 20},
	{FreeCells: solver.MaxFreeCellsForValidBoard, TopN: 20, SelectBest: 20},
}

func hasEnoughFinalCandidates(finalCandidates *internal.SortedBoardStates, requestedCount int) bool {
	if requestedCount <= 0 {
		return finalCandidates.Size() > 0
	}
	return finalCandidates.Size() >= requestedCount
}

func (g *Generator) generateSlow(ctx context.Context) []*solver.Result {
	tries := 0

	if g.count <= 0 {
		g.count = 100
	}
	finalCandidates := internal.NewSortedBoardStates(g.count)
	finalCandidatesReportedToStats := 0
	var stageStats stats.GameStages

	// Start with basic boards in a easy range, so we can generate and sort a lot of
	// candidates. If we start with harder boards, the generation becomes slower.
	start := time.Now()

	// by default, replace the solution every 10 tries, but with the higher levels
	// prefer to keep it longer to benefit more from the cache
	replaceSolutionEvery := 10
	if g.lr.Min >= solver.LevelNightmare {
		replaceSolutionEvery = 50
	} else if g.lr.Min >= solver.LevelDarkEvil {
		replaceSolutionEvery = 25
	}

	// turn on use cache
	initState := g.newInitialBoardState(ctx, true)
	candidates := internal.NewSortedBoardStates(1000)

	for ctx.Err() == nil {
		tries++

		candidates.Reset()
		startState := initState
		startState = g.removeSingleValue(ctx, initState)
		candidates.Add(startState)

		// enhance the candidates to the desired level
		for si := 0; si < len(slowStages) && ctx.Err() == nil; si++ {
			stageStats.ReportCandidateCount(si, candidates.Size())
			if candidates.Size() == 0 {
				// stop after reporting the empty stage
				break
			}
			stage := slowStages[si]
			newFinal, newCandidates, bestComplexity := g.generateSlowStage(ctx, candidates, stage, g.count)
			if bestComplexity > 0 {
				// report the complexity of the candidates per stage
				stageStats.ReportBestComplexity(si, int64(bestComplexity))
			}
			if newFinal.Size() == 0 && (newCandidates.Size() == 0 || si == len(slowStages)-1) {
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
			candidates = newCandidates
			if hasEnoughFinalCandidates(finalCandidates, g.count) {
				break
			}
		}
		// before we give up on this round - we may have in-flight candidates that made the bar
		// but not reported in any stage, do it now
		for bs := range candidates.Boards {
			if !bs.Progress().InRange() {
				// candidates are sorted
				break
			}
			finalCandidates.Add(bs)
			if g.onNewResult != nil {
				g.onNewResult(bs.Result())
			}
		}

		if hasEnoughFinalCandidates(finalCandidates, g.count) {
			break
		}

		if tries >= replaceSolutionEvery {
			// if we are stuck on the same solution for too long, try the next one
			// it also helps reducing the per-solution cache footprint
			cacheStats := initState.SolutionState().Cache().Stats()
			newFinalCandidatesToReport := finalCandidates.Size() - finalCandidatesReportedToStats
			finalCandidatesReportedToStats = finalCandidates.Size()
			stats.Stats.ReportGeneration(newFinalCandidatesToReport,
				time.Since(start), int64(tries), stageStats, cacheStats)
			initState = g.newInitialBoardState(ctx, true)
			tries = 0
			stageStats = stats.GameStages{}
			start = time.Now()
		}
	}

	// return the results so far, even if ctx canceled in the middle
	cacheStats := initState.SolutionState().Cache().Stats()
	newFinalCandidatesToReport := finalCandidates.Size() - finalCandidatesReportedToStats
	stats.Stats.ReportGeneration(newFinalCandidatesToReport, time.Since(start), int64(tries), stageStats, cacheStats)
	return finalCandidates.Results()
}

func (g *Generator) generateSlowStage(
	ctx context.Context,
	candidates *internal.SortedBoardStates,
	stage slowStage,
	maxReady int,
) (
	ready *internal.SortedBoardStates,
	next *internal.SortedBoardStates,
	bestComplexity solver.StepComplexity,
) {
	if stage.SelectBest <= 0 {
		panic("slow stage must have SelectBest > 0")
	}
	if stage.TopN > 0 {
		// using TopN
		topN := internal.TopN(ctx, &internal.TopNArgs{
			In:         candidates,
			FreeCells:  stage.FreeCells,
			TopN:       stage.TopN,
			SelectBest: stage.SelectBest,
		})
		bestComplexity = updateBestComplexityFromBest(bestComplexity, topN.Ready)
		bestComplexity = updateBestComplexityFromBest(bestComplexity, topN.Next)
		return topN.Ready, topN.Next, bestComplexity
	}

	next = internal.NewSortedBoardStates(stage.SelectBest)
	ready = internal.NewSortedBoardStates(maxReady)
	// refine the candidates to the desired level
	for bs := range candidates.Boards {
		switch bs.Candidates().Size() {
		case 0:
			// this board can no longer be enhanced
			if bs.Progress().InRange() {
				ready.Add(bs)
				bestComplexity = updateBestComplexity(bestComplexity, bs)
			}
			// throw this candidate away, it is below the desired level
			continue
		case 1:
			// last candidate, forking is useless - just try to remove it
			bs = bs.RemoveOneByOne(ctx, stage.FreeCells)
			if bs != nil {
				if bs.Progress().InRange() {
					ready.Add(bs)
				}
				// even if we about to throw this candidate away, we still want to report
				// its best complexity
				bestComplexity = updateBestComplexity(bestComplexity, bs)
			}
			continue
		}
		for range stage.GeneratePerCandidate {
			if bs.Candidates() == indexes.MinBitSet81 {
				// no more candidates to remove, we can stop
				bestComplexity = updateBestComplexity(bestComplexity, bs)
				if bs.Progress().InRange() {
					ready.Add(bs)
				}
				break
			}
			if bs.Progress() == internal.InRangeStop {
				// we can stop with this candidate if it is already at the desired level
				ready.Add(bs)
				bestComplexity = updateBestComplexity(bestComplexity, bs)
				break
			}
			var bsForked *internal.BoardState
			// can refine more, try random first
			//if stage.MaxToRemove > 1 {
			bsForked = bs.Remove(ctx, internal.RemoveArgs{
				FreeCells:        stage.FreeCells,
				BatchMinToRemove: stage.MinToRemove,
				BatchMaxToRemove: stage.MaxToRemove,
				BatchMaxTries:    5,
			})
			if bsForked == nil {
				bsForked = bs.RemoveOneByOne(ctx, stage.FreeCells)
			}
			if bsForked == nil {
				// try another removal combination of the same board
				continue
			}
			bestComplexity = updateBestComplexity(bestComplexity, bsForked)
			if bsForked.Progress() == internal.InRangeStop {
				// capture if we have reached the desired level, or if we are at level but
				// also the last stage
				ready.Add(bsForked)
			} else {
				next.Add(bsForked)
			}
		}
	}
	return
}

func updateBestComplexity(bestComplexity solver.StepComplexity, bs *internal.BoardState) solver.StepComplexity {
	if bs == nil {
		return bestComplexity
	}
	return max(bestComplexity, bs.Complexity())
}

func updateBestComplexityFromBest(bestComplexity solver.StepComplexity, sbs *internal.SortedBoardStates) solver.StepComplexity {
	if sbs == nil || sbs.Size() == 0 {
		return bestComplexity
	}
	return updateBestComplexity(bestComplexity, sbs.Get(0))
}
