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

func (g *Generator) generateSlow(ctx context.Context, initState *internal.BoardState, count int) []*solver.Result {
	tries := 0

	candidates := internal.NewSortedBoardStates()
	finalCandidates := internal.NewSortedBoardStates()
	var stageStats internal.GamePerStageStats

	// Start with basic boards in a easy range, so we can generate and sort a lot of
	// candidates. If we start with harder boards, the generation becomes slower.
	start := time.Now()

generationLoop:
	for ctx.Err() == nil {
		candidates.Reset()
		candidates.Add(initState)
		tries++

		// enhance the candidates to the desired level
		for stage := 0; stage < len(slowStages) && ctx.Err() == nil; stage++ {
			newFinal, newCandidates := g.generateSlowStage(ctx, candidates, slowStages[stage])
			if newFinal.Size() == 0 && len(newCandidates) == 0 {
				stageStats.Report(0, stage)
				break
			}
			finalCandidates.AddAll(newFinal)
			if newFinal.Size() > 0 {
				stageStats.Report(newFinal.Size(), stage)
			}
			if hasEnoughFinalCandidates(finalCandidates, count) {
				// stop once we have at least one if count was requested
				internal.Stats.ReportGeneration(finalCandidates.Size(), time.Since(start), int64(tries), stageStats)
				break generationLoop
			}
			candidates = combineCandidates(newCandidates, slowStages[stage].SelectBest)
		}
	}

	// return the results so far, even if ctx canceled in the middle
	internal.Stats.ReportGeneration(finalCandidates.Size(), time.Since(start), int64(tries), stageStats)
	if count > 0 {
		finalCandidates.TrimSize(count)
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
				finalCandidates.Add(bs)
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

/*
type phase struct {
	freeCount            int
	minComplexity        solver.StepComplexity
	enforceMinComplexity bool
	generateCount        int
	selectCount          int
	restore              bool
	merge                bool
}

// TODO: turn off
const debugPrint = true

var phases = []phase{
	// freeCount, minComplexity, enforceMinComplexity, generateCount, selectCount, restore, merge
	{40, 100, true, 1000, 5, true, true},
	{42, 220, true, 20, 5, true, false},
	{44, 450, true, 30, 5, true, false},
	{46, 700, true, 40, 5, true, false},
	{48, 1000, true, 50, 5, true, false},
	{50, 1300, true, 60, 5, true, false},
	{53, 1700, true, 70, 5, true, false},
	{53, 1700, true, 80, 5, true, false},
	{56, 2000, true, 90, 5, true, false},
	{56, 2000, true, 100, 5, true, false},
	// Stop restoration
	{56, 2000, true, 100, 5, false, false},
}

func (g *Generator) generateSlow(ctx context.Context, initialState *internal.BoardState) *solver.Result {
	var current []*internal.BoardState
	for phi, phase := range phases {
		if phi == 0 {
			current = g.generatePhase0Batch(ctx, initialState, phase)
		} else {
			current = g.generatePhaseNBatch(ctx, current, phase)
			if phase.enforceMinComplexity &&
				(phase.minComplexity > current[0].Complexity(ctx)) {
				// TODO: if debugPrint {
				//	printBest(phi, "failed", phase.minComplexity, current)
				//	println("----")
				//}
				// do not bother if requirements were not met.
				return current[0].Solve(ctx)
			}
		}
		// TODO: printBest(phi, "base", phase.minComplexity, current)
		if phase.restore {
			g.restore(ctx, current, phase.minComplexity)
			// TODO: printBest(phi, "restored", phase.minComplexity, current)
		}

		if phase.merge {
			var mergedAtLeastOne bool
			current, mergedAtLeastOne := g.addMergedBoards(ctx, current)
			// TODO: printBest(phi, "merged(before sort)", phase.minComplexity, current);
			if mergedAtLeastOne {
				current = g.sortByComplexityAndTrim(ctx, current, phase.selectCount, phase.minComplexity)
				// TODO: printBest(phi, "merged", phase.minComplexity, current);
				if phase.restore {
					g.restore(ctx, current, phase.minComplexity)
					// TODO: printBest(phi, "merged/restored", phase.minComplexity, current);
				}
			}
		}
	}

	{
		var last []*internal.BoardState
		for _, bLast := range current {
			last = g.findAllPossibleBoards(ctx, bLast, last, 1)
		}
		last = g.sortByComplexityAndTrim(ctx, last, 3, 0)
		current = last
	}
	// TODO: printBest(phase_count, "last", 0, current)
	for _, bs := range current {
		res := bs.Solve(ctx)
		if res.Steps.Level < solver.LevelDarkEvil {
			// current is sorted by level, no point in continuing
			break
		}
		if debugPrint {
			boardSerialized := boards.Serialize(res.Input)
			fmt.Printf("Level: %s, Complexity: %d, Board: %s\n",
				res.Steps.Level, res.Steps.Complexity, boardSerialized)
			// TODO: do we need generation chain here?
		}
	}
	return current[0].Solve(ctx)
}

func (g *Generator) findAllPossibleBoards(
	ctx context.Context, source *internal.BoardState,
	best []*internal.BoardState, depth int,
) []*internal.BoardState {
	mustStayIndexes := make([]int, 0)
	type candidate struct {
		bs    *internal.BoardState
		index int
	}
	var candidates []candidate

	for _, index := range source.RemainedIndexes() {
		state := source.Clone()
		if state.TryRemoveFromIndex(ctx, index) {
			candidates = append(candidates, candidate{state, index})
		} else {
			mustStayIndexes = append(mustStayIndexes, index)
		}
	}

	if len(candidates) == 0 {
		// Leaf board - stop now.
		best = append(best, source)
		return best
	}

	maxCandidates := getMaxCandidateSize(depth)

	if len(candidates) > maxCandidates {
		// recursion depth will explode, limit it.
		internal.RandShuffle(g.r, candidates)
		candidates = candidates[:maxCandidates]
	}

	// We have left with states with valid boards, try to fine tune one-by-one only those boards, removing
	// indexes previously found as "must stay".
	for _, candidate := range candidates {
		for _, index := range mustStayIndexes {
			candidate.bs.RemoveIndex(index)
		}

		// recursively find more boards.
		best = g.findAllPossibleBoards(ctx, candidate.bs, best, depth+1)
		mustStayIndexes = append(mustStayIndexes, candidate.index)
	}
	// once in a while, get rid of duplicates
	if depth%3 == 0 {
		best = g.sortByComplexityAndTrim(ctx, best, len(best), 0)
	}
	return best
}

func (g *Generator) restore(ctx context.Context, best []*internal.BoardState, minComplexity solver.StepComplexity) {
	for i, bs := range best {
		bs = bs.Clone()
		bs.RestoreSimpleValues(ctx, minComplexity)
		best[i] = bs
	}
}

func (g *Generator) sortByComplexityAndTrim(ctx context.Context, states []*internal.BoardState, trimTo int, minComplexity solver.StepComplexity) []*internal.BoardState {
	slices.SortFunc(states, func(state1, state2 *internal.BoardState) int {
		c1 := state1.Complexity(ctx)
		c2 := state2.Complexity(ctx)
		return int(c2 - c1) // sort descending
	})
	// get rid of duplicates
	{
		uniqueStates := make([]*internal.BoardState, 0, len(states))
		var compareTo []*internal.BoardState
		for _, state := range states {
			if len(compareTo) == 0 || state.Complexity(ctx) != compareTo[0].Complexity(ctx) {
				// states are sorted by complexity, so if we have a new complexity, we can reset the whole compareTo
				compareTo = compareTo[:0]
				compareTo = append(compareTo, state)
				uniqueStates = append(uniqueStates, state)
				continue
			}

			// check if we have a dupe
			var foundDupe bool
			for _, ctState := range compareTo {
				if state.BoardEquivalentTo(ctState) {
					foundDupe = true
					break
				}
			}
			if foundDupe {
				// skip this state, it is a duplicate
				continue
			}
			// not a dupe, add to compareTo and uniqueStates
			compareTo = append(compareTo, state)
			uniqueStates = append(uniqueStates, state)
		}
		states = uniqueStates
	}

	trimTo = min(trimTo, len(states))
	for trimTo > 1 &&
		minComplexity > states[trimTo-1].Complexity(ctx) {
		trimTo--
	}
	states = states[:trimTo]
	return states
}

func (g *Generator) addMergedBoards(
	ctx context.Context, states []*internal.BoardState) ([]*internal.BoardState, bool) {
	// merge boards with the same complexity
	mergedAtLeastOne := false
	// capture size before adding new ones.
	originalSize := len(states)
	for i := 0; (i + 1) < originalSize; i++ {
		for j := i + 1; j < originalSize; j++ {
			mergedState := internal.TryMergeBoardStates(ctx, states[i], states[j])
			if mergedState != nil {
				// we have a new state, add it to the list
				states = append(states, mergedState)
				mergedAtLeastOne = true
			}
		}
	}

	return states, mergedAtLeastOne
}

func getMaxCandidateSize(depth int) int {
	// 5 * 5 * 5 * 4 * 3 * 2 = 3,000
	if depth < 4 {
		return 5
	}
	if depth < 7 {
		return 8 - depth
	}

	return 1
}

func (g *Generator) getBestWithRetries(ctx context.Context,
	best []*internal.BoardState, start *internal.BoardState, count int,
	args internal.RemoveWithRetriesArgs) []*internal.BoardState {
	// ensure start has a valid last result
	start.Solve(ctx)

	newAdded := 0
	betterFound := false

	for newAdded < count {
		state := start.Clone()
		progress := state.RemoveWithRetries(ctx, args)
		if progress.KeepGoing() || progress == internal.AtLevelStop {
			break
		}

		newAdded++
		best = append(best, state)
		betterFound = betterFound || state.Complexity(ctx) > start.Complexity(ctx)
	}

	if !betterFound {
		best = append(best, start.Clone())
	}

	return best

}

func (g *Generator) generatePhase0Batch(
	ctx context.Context,
	start *internal.BoardState,
	phase phase) []*internal.BoardState {
	best := make([]*internal.BoardState, 0, phase.generateCount)
	done := false
	count := 0
	for !done {
		state := g.removePhase0(ctx, start, phase.freeCount)
		best = append(best, state)
		count++
		if count < phase.generateCount {
			done = false
		} else {
			best = g.sortByComplexityAndTrim(ctx, best, phase.selectCount, phase.minComplexity)
			if len(best) == 0 {
				done = false
			} else {
				done = !phase.enforceMinComplexity ||
					phase.minComplexity <= best[0].Complexity(ctx)
			}
		}
	}
	return best
}

func (g *Generator) removePhase0(ctx context.Context, start *internal.BoardState, freeCount int) *internal.BoardState {
	for {
		state := start.Clone()
		// Tests show that the fastest and the most efficient way to remove the first batch appears to be
		// a [2/3] combination.
		progress := state.RemoveWithRetries(ctx, internal.RemoveWithRetriesArgs{
			FreeAtLeast: freeCount - 2,
			MinToRemove: 2,
			MaxToRemove: 3,
			MaxRetries:  15,
		})
		if progress.KeepGoing() {
			progress = state.RemoveWithRetries(ctx, internal.RemoveWithRetriesArgs{
				FreeAtLeast: freeCount,
				MinToRemove: 1,
				MaxToRemove: 1,
				MaxRetries:  15,
			})
		}
		if progress.KeepGoing() || progress == internal.AtLevelStop {
			// done with phase 0
			return state
		}
	}
}

func (g *Generator) generatePhaseNBatch(
	ctx context.Context,
	sources []*internal.BoardState,
	phase phase) []*internal.BoardState {
	next := make([]*internal.BoardState, 0, phase.generateCount)
	for _, bs := range sources {
		maxRetries := min(30, bs.RemainedIndexesSize())
		bs.Clone()
		minToRemove := 1
		maxToRemove := 2
		if bs.RemainedIndexesSize() > 44 {
			minToRemove++
			maxToRemove++
		}
		next = g.getBestWithRetries(
			ctx, next, bs,
			phase.generateCount, internal.RemoveWithRetriesArgs{
				FreeAtLeast: phase.freeCount,
				MinToRemove: minToRemove,
				MaxToRemove: maxToRemove,
				MaxRetries:  maxRetries,
			})
	}
	next = g.sortByComplexityAndTrim(ctx, next, phase.selectCount, phase.minComplexity)
	return next
}
*/
