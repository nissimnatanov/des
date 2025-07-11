package internal

import (
	"context"
	"fmt"
	"slices"

	"github.com/nissimnatanov/des/go/internal/collections"
	"github.com/nissimnatanov/des/go/solver"
)

// SortedBoards are sorted by complexity and deduplicated
type SortedBoardStates struct {
	maxSize int
	sorted  []*BoardState
}

func NewSortedBoardStates(maxSize int) *SortedBoardStates {
	if maxSize <= 0 {
		panic("maxSize must be greater than 0")
	}
	return &SortedBoardStates{
		sorted:  make([]*BoardState, 0, maxSize),
		maxSize: maxSize,
	}
}

func (sbs *SortedBoardStates) Reset() {
	sbs.sorted = sbs.sorted[:0]
}

func (sbs *SortedBoardStates) Size() int {
	return len(sbs.sorted)
}

func (sbs *SortedBoardStates) Get(index int) *BoardState {
	return sbs.sorted[index]
}

func (sbs *SortedBoardStates) Boards(yield func(*BoardState) bool) {
	for _, bs := range sbs.sorted {
		if !yield(bs) {
			return
		}
	}
}

func (sbs *SortedBoardStates) SolveAll(ctx context.Context) {
	if len(sbs.sorted) == 0 || sbs.sorted[0].Action() == solver.ActionSolve {
		return
	}
	for _, bs := range sbs.sorted {
		bs.Solve(ctx)
	}
	// before we sort, throw away the overflow boards
	sbs.sorted = slices.DeleteFunc(sbs.sorted, func(bs *BoardState) bool {
		return bs.Progress() == AboveMaxLevel
	})

	// now re-sort the boards
	slices.SortFunc(sbs.sorted, func(a, b *BoardState) int {
		return int(b.Complexity() - a.Complexity())
	})
}

// Add adds a new board to the sorted list if it is not a duplicate.
func (sbs *SortedBoardStates) Add(newBoard *BoardState) bool {
	_, added := sbs.addInternal(newBoard, 0)
	return added
}

func (sbs *SortedBoardStates) addInternal(newBoard *BoardState, from int) (int, bool) {
	if len(sbs.sorted) > 0 && sbs.sorted[0].Action() != newBoard.Action() {
		// if the first board is below the current level, we cannot add any new boards
		panic(fmt.Sprintf(
			"cannot add a board with a different action to the sorted list: had %s, wanted %s",
			sbs.sorted[0].Action(), newBoard.Action()))
	}
	if len(sbs.sorted) == sbs.maxSize {
		if sbs.sorted[sbs.maxSize-1].Complexity() >= newBoard.Complexity() {
			// do not add boards that are below the current min complexity if at capacity
			return -1, false
		}
		// at capacity, adding new boards will require removing the last one
		sbs.sorted = sbs.sorted[:sbs.maxSize-1]
	}

	insertIndex := len(sbs.sorted)
sortingLoop:
	for i := from; i < len(sbs.sorted); i++ {
		bs := sbs.sorted[i]
		switch {
		case bs.Complexity() > newBoard.Complexity():
			continue
		case bs.Complexity() < newBoard.Complexity():
			insertIndex = i
			break sortingLoop
		}
		// do not add duplicates
		if bs.BoardEquivalentTo(newBoard) {
			return i, false
		}
	}
	sbs.sorted = slices.Insert(sbs.sorted, insertIndex, newBoard)
	return insertIndex, true
}

func (sbs *SortedBoardStates) AddAll(boards *SortedBoardStates) {
	lastIndex := -1
	for _, board := range boards.sorted {
		lastIndex, _ = sbs.addInternal(board, lastIndex+1)
	}
}

func (sbs *SortedBoardStates) Results() []*solver.Result {
	return collections.MapSlice(sbs.sorted, (*BoardState).Result)
}
