package internal

import (
	"context"
	"fmt"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/solver"
)

type BoardState struct {
	board    *boards.Game
	state    *State
	res      *solver.Result
	progress Progress

	// candidates are the board indices with values that can be potentially removed
	// to reach the desired level
	candidates indexes.BitSet81
}

// newBoardState creates a new BoardState for the given board.
func newBoardState(ctx context.Context, state *State, board *boards.Game) *BoardState {
	res := state.solver.Run(ctx, board)
	if res.Status != solver.StatusSucceeded {
		panic(fmt.Sprintf("failed to solve the board: %s", res.Error))
	}
	bs := &BoardState{
		board:      board,
		state:      state,
		res:        res,
		progress:   shouldContinue(state, board, res),
		candidates: board.EmptyCells().Complement(),
	}
	bs.checkIntegrity()
	return bs
}

func (bs *BoardState) checkIntegrity() {
	if !boards.GetIntegrityChecks() {
		return
	}

	allIndexes := make([]int, 0, boards.Size)
	for index := range bs.candidates.Indexes {
		if bs.board.Get(index) == 0 {
			panic("remained index points to an empty value")
		}
		allIndexes = append(allIndexes, index)
	}
}

func (bs *BoardState) Complexity(ctx context.Context) solver.StepComplexity {
	return bs.res.Steps.Complexity
}

func (bs *BoardState) Level(ctx context.Context) solver.Level {
	return bs.res.Steps.Level
}

func (bs *BoardState) Result() *solver.Result {
	return bs.res
}

func (bs *BoardState) Progress() Progress {
	return bs.progress
}

func (bs *BoardState) BoardEquivalentTo(other *BoardState) bool {
	return boards.Equivalent(bs.board, other.board)
}

func (bs *BoardState) Candidates() indexes.BitSet81 {
	return bs.candidates
}

func shouldContinueAtLevel(desiredLevel solver.Level, r *Random) bool {
	switch desiredLevel {
	case solver.LevelEasy:
		// For easy games - keep trying (otherwise, game can be too easy).
		return r.PercentProbability(95)
	case solver.LevelMedium:
		// For medium games - keep trying a bit less.
		return r.PercentProbability(75)
	case solver.LevelHard:
		// For hard games - continue in half of the cases..
		return r.PercentProbability(50)
	case solver.LevelVeryHard:
		// For very hard games - make it even harder, but stop sometimes.
		return r.PercentProbability(75)
	default:
		// For harder games, keep going until overflows...
		return true
	}
}

func shouldContinue(state *State, board *boards.Game, res *solver.Result) Progress {
	if board.FreeCellCount() < 32 {
		// too early even for easy games.
		return TooEarly
	}

	if res.Steps.Level < state.level {
		return BelowLevel
	}

	if res.Steps.Level > state.level {
		// Overflow, stop.
		return AboveLevel
	}

	if shouldContinueAtLevel(state.level, state.rand) {
		// Keep going, we are at the desired level.
		return AtLevelKeepGoing
	}

	// We are at the desired level, but do not want to continue.
	return AtLevelStop
}

func (bs *BoardState) Remove(ctx context.Context, args RemoveArgs) *BoardState {
	if bs.progress == AtLevelStop || bs.progress == AboveLevel {
		// we already overflowed the level, no point in removing anything
		panic("do not use Remove if already reached the desired level or overflowed it")
	}
	defer bs.checkIntegrity()

	next := bs
	removedOnce := false
	for next.board.FreeCellCount() < args.FreeAtLeast {
		nextRemoved := next.tryRemove(ctx, &args)
		if nextRemoved == nil {
			break
		}
		removedOnce = true
		next = nextRemoved
		// make sure we do not overflow the level
		switch next.progress {
		case TooEarly, BelowLevel, AtLevelKeepGoing:
			// keep removing while we reach the desired level or FreeAtLeast threshold
			continue
		case AtLevelStop:
			// we reached the desired level, stop removing even if we have not reached the FreeAtLeast
			return next
		case AboveLevel:
			panic("tryRemove should not overflow the level")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemove: %d", next.progress))
		}
	}
	if !removedOnce {
		// even if we did not remove anything, we may have reached the desired level
		if next.progress == AtLevelKeepGoing {
			// we reached the desired level, but did not remove anything
			// let the caller know it is time to stop
			next.progress = AtLevelStop
			return next
		}
		return nil
	}

	return next
}

// RemoveOneByOne tries to remove indexes one by one, until the board reaches the desired level
// or the number of free cells is less than MaxFreeCellsForValidBoard.
func (bs *BoardState) RemoveOneByOne(ctx context.Context) *BoardState {
	if bs.progress == AtLevelStop || bs.progress == AboveLevel {
		// we already overflowed the level, no point in removing anything
		panic("do not use RemoveOneByOne if already reached the desired level or overflowed it")
	}
	defer bs.checkIntegrity()

	// Remove the remained indexes one by one, until we reach the desired level.
	r := bs.state.rand
	candidates := slices.Collect(bs.candidates.Indexes)
	if len(candidates) == 0 {
		if bs.progress == AtLevelKeepGoing {
			bs.progress = AtLevelStop
			return bs
		}
		return nil
	}

	RandShuffle(r, candidates)
	removedOnce := false
	next := bs
	for ci := range candidates {
		next2 := next.tryRemoveCandidates(ctx, candidates[ci:ci+1])
		if next2 == nil {
			continue
		}
		next = next2
		removedOnce = true
		// tryRemoveOne does not overflow the level
		switch next.progress {
		case TooEarly, BelowLevel, AtLevelKeepGoing:
			// keep removing more
		case AtLevelStop:
			// we are done
			return next
		case AboveLevel:
			panic("tryRemoveOne should not overflow the level")
		default:
			panic(fmt.Sprintf("unexpected progress value after tryRemoveOne: %d", next.progress))
		}
	}
	// if we tried all the candidates and reached the level, we can stop
	if next.progress == AtLevelKeepGoing {
		next.progress = AtLevelStop
		return next
	}
	if !removedOnce {
		// return what we have so far
		return nil
	}

	return next
}

func (bs *BoardState) tryRemove(ctx context.Context, args *RemoveArgs) *BoardState {
	if args.MinToRemove < 1 || args.MaxToRemove < args.MinToRemove {
		panic("minToRemove and maxToRemove are out of range")
	}
	defer bs.checkIntegrity()
	r := bs.state.rand

	next := bs
	removedOnce := false
	for range args.MaxRetries {
		allowedToRemove := solver.MaxFreeCellsForValidBoard - next.board.FreeCellCount()
		if allowedToRemove == 0 {
			break
		}
		candidates := slices.Collect(next.candidates.Indexes)
		if len(candidates) == 0 {
			break
		}

		RandShuffle(r, candidates)
		currentBatch := r.NextInClosedRange(args.MinToRemove, args.MaxToRemove)
		currentBatch = min(currentBatch, len(candidates))
		currentBatch = min(currentBatch, allowedToRemove)

		nextRemoved := next.tryRemoveCandidates(ctx, candidates[:currentBatch])
		if nextRemoved == nil {
			continue // try again
		}
		next = nextRemoved
		removedOnce = true
		switch next.progress {
		case TooEarly, BelowLevel, AtLevelKeepGoing, AtLevelStop:
			// we removed, all good, caller will call again
			return next
		case AboveLevel:
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
		res := bs.state.prover.Run(ctx, bs.board)
		if res.Status != solver.StatusSucceeded {
			panic("do not use invalid boards as an input here")
		}
	}

	for _, index := range candidates {
		bs.board.Set(index, 0)
	}

	res := bs.state.prover.Run(ctx, bs.board)
	var nextBoard *boards.Game
	if res.Status == solver.StatusSucceeded {
		// clone the new board
		nextBoard = bs.board.Clone(boards.Edit)
	}
	// always restore the board to its original state
	for _, index := range candidates {
		bs.board.SetReadOnly(index, bs.state.solution.Get(index))
	}
	var next *BoardState
	if nextBoard != nil {
		next = newBoardState(ctx, bs.state, nextBoard)
		// if we overflowed the desired level, consider it a failure as well
		if next.progress == AboveLevel {
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
