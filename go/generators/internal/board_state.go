package internal

import (
	"context"
	"fmt"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/solver"
)

type BoardState struct {
	solState *SolutionState
	res      *solver.Result

	desiredLevelRange LevelRange
	// progress is relative to the desired level range
	progress Progress

	// candidates are the board indices with values that can be potentially removed
	// to reach the desired level range
	candidates indexes.BitSet81
}

func (bs *BoardState) tryCleanIndexes(ctx context.Context, indexesToClean []int, proveOnly bool) *BoardState {
	if len(indexesToClean) == 0 {
		panic("indexesToClean cannot be empty")
	}
	defer bs.checkIntegrity()

	b := bs.board()
	for _, index := range indexesToClean {
		b.Set(index, 0)
	}

	// we could prob create fake result for solutions, but it does not matter much
	// if we know current board required at least some recursion depth,
	// we can pass it to the solver of the child board to save tons of time on redundant sub-layer tries
	action := solver.ActionSolve
	minRecDepth := int8(0)
	if bs.res.Action == solver.ActionSolve {
		// we can only use the recursion depth from the previous result if it was a Solve action
		// Prove does not have limits and its recursion depth is not relevant and can be too deep
		minRecDepth = bs.res.RecursionDepth
	}
	if proveOnly {
		action = solver.ActionProve
	}
	res := bs.solState.solver.Run(ctx, b, action, bs.solState.cache,
		solver.WithMinRecursionDepth(minRecDepth))

	// always restore the board to its original state
	for _, index := range indexesToClean {
		b.SetReadOnly(index, bs.solState.solution.Get(index))
	}

	// capture all, including unsolved (those also contribute to the slow generation)
	SlowBoards.Add(res)

	if res.Status != solver.StatusSucceeded {
		return nil
	}

	// preserve the clone as Edit board instead of the original board
	editBoard := res.Input.Clone(boards.Edit)
	res.Input = editBoard
	progress := TooEarly
	if action == solver.ActionSolve {
		progress = bs.desiredLevelRange.shouldContinue(bs.solState.rand, editBoard, res)
	}
	// if the parent board had some candidates removed, we should remove them from the child state as well
	childCandidates := editBoard.EmptyCells().Complement().Intersect(bs.candidates)

	child := &BoardState{
		solState:          bs.solState,
		res:               res,
		desiredLevelRange: bs.desiredLevelRange,
		progress:          progress,
		candidates:        childCandidates,
	}
	child.checkIntegrity()
	return child
}

func newSolutionBoardState(
	ctx context.Context, state *SolutionState, levelRange LevelRange, sol *boards.Solution,
) *BoardState {
	// we could prob create fake result for solutions, but it does not matter much
	editBoard := sol.Clone(boards.Edit)
	res := state.solver.Run(ctx, editBoard, solver.ActionProve, state.cache)
	if res.Status != solver.StatusSucceeded {
		panic("failed to solve a solution")
	}
	// flip back to the Edit board
	res.Input = editBoard
	bs := &BoardState{
		solState:          state,
		res:               res,
		desiredLevelRange: levelRange,
		progress:          TooEarly,
		candidates:        indexes.MaxBitSet81,
	}
	bs.checkIntegrity()
	return bs
}

func (bs *BoardState) board() *boards.Game {
	return bs.res.Input
}

func (bs *BoardState) checkIntegrity() {
	if !boards.GetIntegrityChecks() {
		return
	}

	allIndexes := make([]int, 0, boards.Size)
	for index := range bs.candidates.Indexes {
		if bs.board().Get(index) == 0 {
			panic("remained index points to an empty value")
		}
		allIndexes = append(allIndexes, index)
	}
	if !boards.ContainsAll(bs.solState.solution, bs.board()) {
		panic("provided solution does not contain the board")
	}
}

func (bs *BoardState) SolutionState() *SolutionState {
	return bs.solState
}

func (bs *BoardState) Action() solver.Action {
	return bs.res.Action
}

func (bs *BoardState) Complexity() solver.StepComplexity {
	return bs.res.Complexity
}

func (bs *BoardState) Result() *solver.Result {
	return bs.res
}

func (bs *BoardState) Progress() Progress {
	return bs.progress
}

func (bs *BoardState) BoardLevel() solver.Level {
	return bs.res.Level
}

func (bs *BoardState) DesiredLevelRange() LevelRange {
	return bs.desiredLevelRange
}

func (bs *BoardState) BoardEquivalentTo(other *BoardState) bool {
	return boards.Equivalent(bs.board(), other.board())
}

func (bs *BoardState) Candidates() indexes.BitSet81 {
	return bs.candidates
}

type RemoveArgs struct {
	FreeCells        int
	BatchMinToRemove int
	BatchMaxToRemove int
	BatchMaxTries    int
	ProveOnly        bool
}

func (bs *BoardState) Remove(ctx context.Context, args RemoveArgs) *BoardState {
	if bs.progress == AboveMaxLevel {
		// we already overflowed the level, no point in removing anything
		panic("do not use Remove if the board is already above the max level")
	}
	if bs.progress == InRangeStop && bs.candidates == indexes.MinBitSet81 {
		// we can no longer tune the board and it is on desired level, hence
		// calling this method in a loop will cause an infinite loop
		panic("do not use Remove if already reached the desired level or overflowed it")
	}
	defer bs.checkIntegrity()

	next := bs
	removedOnce := false
	for next.board().FreeCellCount() < args.FreeCells {
		{
			nextRemoved := next.tryRemove(ctx, &args)
			if nextRemoved == nil {
				break
			}
			removedOnce = true
			next = nextRemoved
		}
		// make sure we do not overflow the level
		switch next.progress {
		case TooEarly, BelowMinLevel, InRangeKeepGoing:
			// keep removing while we reach the desired level or FreeAtLeast threshold
			continue
		case InRangeStop:
			// we reached the desired level, stop removing even if we have not reached the FreeAtLeast
			return next
		case AboveMaxLevel:
			panic("tryRemove should not overflow the level")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemove: %d", next.progress))
		}
	}
	if !removedOnce {
		// even if we did not remove anything, we may have reached the desired level
		if next.progress == InRangeKeepGoing {
			// we reached the desired level, but did not remove anything
			// let the caller know it is time to stop
			next.progress = InRangeStop
			return next
		}
		return nil
	}

	return next
}

// RemoveOneByOne tries to remove indexes one by one, until the board reaches the desired level
// or the number of free cells is less than MaxFreeCellsForValidBoard.
func (bs *BoardState) RemoveOneByOne(ctx context.Context, freeCells int, proveOnly bool) *BoardState {
	if bs.progress == InRangeStop || bs.progress == AboveMaxLevel {
		// we already overflowed the level, no point in removing anything
		panic("do not use RemoveOneByOne if already reached the desired level or overflowed it")
	}

	// Remove the remained indexes one by one, until we reach the desired level.
	candidates := slices.Collect(bs.candidates.Indexes)
	if len(candidates) == 0 {
		if bs.progress == InRangeKeepGoing {
			bs.progress = InRangeStop
			return bs
		}
		return nil
	}

	defer bs.checkIntegrity()

	r := bs.solState.rand
	random.Shuffle(r, candidates)
	next := bs
	for ci := range candidates {
		if next.board().FreeCellCount() >= freeCells {
			return bs
		}

		{
			nextRemoved := next.tryRemoveCandidates(ctx, candidates[ci:ci+1], proveOnly)
			if nextRemoved == nil {
				continue
			}
			next = nextRemoved
		}
		// tryRemoveOne does not overflow the level
		done := false
		switch next.progress {
		case TooEarly, BelowMinLevel, InRangeKeepGoing:
			// keep removing more
		case InRangeStop:
			// we are done
			done = true
		case AboveMaxLevel:
			panic("tryRemoveOne should not overflow the level")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemoveOne: %d", next.progress))
		}
		if done {
			break
		}
	}
	// if we tried all the candidates and reached the level, we can stop
	if next.progress.InRange() {
		next.progress = InRangeStop
		return next
	}

	// no more candidates to remove and and we did not reach the desired level
	return nil
}

func (bs *BoardState) tryRemove(ctx context.Context, args *RemoveArgs) *BoardState {
	if args.BatchMinToRemove < 1 || args.BatchMaxToRemove < args.BatchMinToRemove {
		panic("minToRemove and maxToRemove are out of range")
	}
	defer bs.checkIntegrity()
	r := bs.solState.rand

	next := bs
	removedOnce := false
	for range args.BatchMaxTries {
		allowedToRemove := solver.MaxFreeCellsForValidBoard - next.board().FreeCellCount()
		if allowedToRemove == 0 {
			break
		}
		candidates := slices.Collect(next.candidates.Indexes)
		if len(candidates) == 0 {
			break
		}

		random.Shuffle(r, candidates)
		currentBatch := r.NextInClosedRange(args.BatchMinToRemove, args.BatchMaxToRemove)
		currentBatch = min(currentBatch, len(candidates))
		currentBatch = min(currentBatch, allowedToRemove)

		{
			nextRemoved := next.tryRemoveCandidates(ctx, candidates[:currentBatch], args.ProveOnly)
			if nextRemoved == nil {
				continue // try again
			}
			removedOnce = true
			next = nextRemoved
		}
		switch next.progress {
		case TooEarly, BelowMinLevel, InRangeKeepGoing, InRangeStop:
			// we removed, all good, caller will call again
			return next
		case AboveMaxLevel:
			panic("tryRemoveBatch should not overflow")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemoveBatch: %d", next.progress))
		}
	}
	if !removedOnce {
		// let the caller know that we did not remove anything
		return nil
	}

	return next
}

func (bs *BoardState) Solve(ctx context.Context) {
	defer bs.checkIntegrity()
	// capture the original editable board
	editBoard := bs.board()
	if bs.res.Action == solver.ActionSolve {
		return
	}
	res := bs.solState.solver.Run(ctx, bs.board(), solver.ActionSolve, bs.solState.cache)
	// capture all, including unsolved (those also contribute to the slow generation)
	SlowBoards.Add(res)

	if res.Status != solver.StatusSucceeded {
		if ctx.Err() == nil {
			panic("cannot solve the board, which was proven before")
		} else {
			return
		}
	}

	// revert back to the editable board
	res.Input = editBoard
	bs.res = res
	bs.progress = bs.desiredLevelRange.shouldContinue(bs.solState.rand, editBoard, res)
	// do not reset candidates, we did not change the board and candidates are still valid
}

// tryRemoveCandidates tries to remove a batch of indexes from the board, only once.
// If it fails, it returns nil.
// If successful, it returns a new state board
// If level overflows, this method reverts back and returns Failed.
// If batchSize is 1, the index is removed from the remained list.
func (bs *BoardState) tryRemoveCandidates(ctx context.Context, candidates []int, proveOnly bool) *BoardState {
	defer bs.checkIntegrity()

	if boards.GetIntegrityChecks() {
		res := bs.solState.solver.Run(ctx, bs.board(), solver.ActionProve, bs.solState.cache)
		if res.Status != solver.StatusSucceeded {
			panic("do not use invalid boards as an input here")
		}
		if bs.Complexity() != res.Complexity {
			panic(fmt.Sprintf("complexity mismatch: %d != %d, board has already changed",
				bs.Complexity(), res.Complexity))
		}
	}

	next := bs.tryCleanIndexes(ctx, candidates, proveOnly)

	if next != nil {
		// if we overflowed the desired level, consider it a failure as well
		if next.progress == AboveMaxLevel {
			next = nil
		}
	}
	if next != nil {
		// prove succeeded and we did not overflow the level
		return next
	}

	// on failure, let's see if we tried to remove only one cell
	// if yes, it is guaranteed that removal of this cell will not succeed further and
	// we can simply remove it from the candidates list of the current board
	if len(candidates) == 1 {
		bs.candidates.Reset(candidates[0])
	}

	return next
}

func (bs *BoardState) WithDesiredLevelRange(lr LevelRange) *BoardState {
	clone := *bs
	if bs.desiredLevelRange == lr {
		// no need to clone, we are already at the desired level
		return &clone
	}
	// reevaluate the progress based on the new level
	clone.desiredLevelRange = lr
	clone.progress = lr.shouldContinue(bs.solState.rand, bs.board(), bs.res)
	return &clone
}

func (bs *BoardState) RemoveVal(ctx context.Context, v values.Value, count int) *BoardState {
	// RemoveVal is invoked in the beginning of the generation process, we can skip res check
	if bs.board().FreeCellCount() > 0 {
		panic("use RemoveVal only when the board is full")
	}
	vIndexes := make([]int, 0, count)
	for i := range boards.Size {
		if bs.board().Get(i) == v {
			vIndexes = append(vIndexes, i)
		}
	}
	if count < len(vIndexes) {
		random.Shuffle(bs.solState.rand, vIndexes)
		vIndexes = vIndexes[:count]
	}

	return bs.tryCleanIndexes(ctx, vIndexes, true)
}

type TopNArgs struct {
	In         *SortedBoardStates
	TopN       int
	FreeCells  int
	SelectBest int
	ProveOnly  bool
}

type TopNResult struct {
	Next           *SortedBoardStates
	Ready          *SortedBoardStates
	BestComplexity solver.StepComplexity
}

func TopN(ctx context.Context, args *TopNArgs) TopNResult {
	if args.In == nil || args.In.Size() == 0 {
		panic("In must be provided for TopN")
	}
	result := TopNResult{
		Next:  NewSortedBoardStates(args.SelectBest),
		Ready: NewSortedBoardStates(args.SelectBest),
	}

	next := args.In
	nextPerCandidate := NewSortedBoardStates(args.TopN)
	for next.Size() > 0 {
		cur := next
		next = NewSortedBoardStates(args.SelectBest)
		for bi := range cur.Size() {
			bs := cur.Get(bi)
			if bs.Progress() == InRangeKeepGoing && bs.Candidates() == indexes.MinBitSet81 {
				bs.progress = InRangeStop
			}
			switch {
			case bs.progress == AboveMaxLevel:
				continue
			case bs.progress == InRangeStop:
				// reached the level
				result.Ready.Add(bs)
				continue
			case bs.board().FreeCellCount() >= args.FreeCells:
				// reached the desired free cell
				result.Next.Add(bs)
				continue
			}

			// try the remained indexes one by one
			nextPerCandidate.Reset()
			for index := range bs.candidates.Indexes {
				removed := bs.tryRemoveCandidates(ctx, []int{index}, args.ProveOnly)
				if removed == nil {
					continue
				}
				result.BestComplexity = MaxComplexity(result.BestComplexity, removed)
				if removed.progress == InRangeStop {
					// we reached the desired level, return it as a ready board
					result.Ready.Add(removed)
					continue
				}
				nextPerCandidate.Add(removed)
				next.Add(removed)
			}
			if nextPerCandidate.Size() == 0 && bs.progress == InRangeKeepGoing {
				// if the board in range and we could not enhance it further, capture as ready
				bs.progress = InRangeStop
				result.Ready.Add(bs)
			}
			for nbs := range nextPerCandidate.Boards {
				// if we tried to remove the candidate index from parent and it lead to a failure,
				// we should remove it from the child state as well
				nbs.candidates = nbs.candidates.Intersect(bs.candidates)
			}
		}
	}
	return result
}
