package internal

import (
	"context"
	"fmt"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/solver"
)

type BoardState struct {
	state *SolutionState
	res   *solver.Result

	desiredLevelRange LevelRange
	// progress is relative to the desired level range
	progress Progress

	// candidates are the board indices with values that can be potentially removed
	// to reach the desired level range
	candidates indexes.BitSet81
}

// newBoardState creates a new BoardState for the given board.
// if solver fails, it returns nil with the failed result
func newBoardState(
	ctx context.Context, state *SolutionState, levelRange LevelRange, srcBoard *boards.Game,
) (*BoardState, *solver.Result) {
	// we could prob create fake result for solutions, but it does not matter much
	res := state.solver.Run(ctx, srcBoard)
	if res.Status != solver.StatusSucceeded {
		return nil, res
	}
	// preserve the clone as Edit board instead of the original board
	// TODO: maybe reset and reuse the play-board since it is no longer needed
	editBoard := res.Input.Clone(boards.Edit)
	res.Input = editBoard
	progress := levelRange.shouldContinue(state.rand, editBoard, res)
	bs := &BoardState{
		state:             state,
		res:               res,
		desiredLevelRange: levelRange,
		progress:          progress,
		candidates:        editBoard.EmptyCells().Complement(),
	}
	bs.checkIntegrity()
	return bs, res
}

func newSolutionBoardState(
	ctx context.Context, state *SolutionState, levelRange LevelRange, sol *boards.Solution,
) *BoardState {
	// we could prob create fake result for solutions, but it does not matter much
	editBoard := sol.Clone(boards.Edit)
	res := state.solver.Run(ctx, editBoard)
	if res.Status != solver.StatusSucceeded {
		panic("failed to solve a solution")
	}
	bs := &BoardState{
		state:             state,
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
	if !boards.ContainsAll(bs.state.solution, bs.board()) {
		panic("provided solution does not contain the board")
	}
}

func (bs *BoardState) Complexity() solver.StepComplexity {
	return bs.res.Steps.Complexity
}

func (bs *BoardState) Result() *solver.Result {
	return bs.res
}

func (bs *BoardState) Progress() Progress {
	return bs.progress
}

func (bs *BoardState) BoardLevel() solver.Level {
	return bs.res.Steps.Level
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
func (bs *BoardState) RemoveOneByOne(ctx context.Context) *BoardState {
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

	return bs.RemoveCandidatesOneByOne(ctx, candidates)
}

func (bs *BoardState) RemoveCandidatesOneByOne(ctx context.Context, candidates []int) *BoardState {
	defer bs.checkIntegrity()

	r := bs.state.rand
	random.Shuffle(r, candidates)
	removedOnce := false
	next := bs
	for ci := range candidates {
		{
			nextRemoved := next.tryRemoveCandidates(ctx, candidates[ci:ci+1])
			if nextRemoved == nil {
				continue
			}
			removedOnce = true
			next = nextRemoved
		}
		// tryRemoveOne does not overflow the level
		switch next.progress {
		case TooEarly, BelowMinLevel, InRangeKeepGoing:
			// keep removing more
		case InRangeStop:
			// we are done
			return next
		case AboveMaxLevel:
			panic("tryRemoveOne should not overflow the level")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemoveOne: %d", next.progress))
		}
	}
	// if we tried all the candidates and reached the level, we can stop
	if next.progress == InRangeKeepGoing {
		next.progress = InRangeStop
		return next
	}
	if !removedOnce {
		// return what we have so far
		return nil
	}

	return next
}

func (bs *BoardState) tryRemove(ctx context.Context, args *RemoveArgs) *BoardState {
	if args.BatchMinToRemove < 1 || args.BatchMaxToRemove < args.BatchMinToRemove {
		panic("minToRemove and maxToRemove are out of range")
	}
	defer bs.checkIntegrity()
	r := bs.state.rand

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
			nextRemoved := next.tryRemoveCandidates(ctx, candidates[:currentBatch])
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

// tryRemoveCandidates tries to remove a batch of indexes from the board, only once.
// If it fails, it returns nil.
// If successful, it returns a new state board
// If level overflows, this method reverts back and returns Failed.
// If batchSize is 1, the index is removed from the remained list.
func (bs *BoardState) tryRemoveCandidates(ctx context.Context, candidates []int) *BoardState {
	defer bs.checkIntegrity()
	if len(candidates) == 0 {
		panic("candidates cannot be empty")
	}

	if boards.GetIntegrityChecks() {
		res := bs.state.prover.Run(ctx, bs.board())
		if res.Status != solver.StatusSucceeded {
			panic("do not use invalid boards as an input here")
		}
	}

	for _, index := range candidates {
		bs.board().Set(index, 0)
	}

	next, _ := newBoardState(ctx, bs.state, bs.desiredLevelRange, bs.board())
	// always restore the board to its original state
	for _, index := range candidates {
		bs.board().SetReadOnly(index, bs.state.solution.Get(index))
	}
	if next != nil {
		// if we overflowed the desired level, consider it a failure as well
		if next.progress == AboveMaxLevel {
			next = nil
		}
	}
	if next != nil {
		// prove succeeded and we did not overflow the level
		// if previous state had some of its candidates removed (because they lead to a failure),
		// we should remove them from the new state as well.
		next.candidates = next.candidates.Intersect(bs.candidates)
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
	clone.progress = lr.shouldContinue(bs.state.rand, bs.board(), bs.res)
	return &clone
}
