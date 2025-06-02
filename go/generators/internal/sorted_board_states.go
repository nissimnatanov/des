package internal

import (
	"slices"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/internal/collections"
	"github.com/nissimnatanov/des/go/solver"
)

// SortedBoards are sorted by complexity and deduplicated
type SortedBoardStates struct {
	sorted []*BoardState
}

func NewSortedBoardStates(bs ...*BoardState) *SortedBoardStates {
	// TODO: can optimize
	sbs := &SortedBoardStates{}
	for _, b := range bs {
		sbs.Add(b)
	}
	return sbs
}

func (sbs *SortedBoardStates) Reset() {
	sbs.sorted = sbs.sorted[:0]
}

func (sbs *SortedBoardStates) Size() int {
	return len(sbs.sorted)
}

func (sbs *SortedBoardStates) TrimSize(limit int) {
	if len(sbs.sorted) < limit {
		return
	}
	sbs.sorted = slices.Delete(sbs.sorted, limit, len(sbs.sorted))
}

func (sbs *SortedBoardStates) Boards(yield func(*BoardState) bool) {
	for _, bs := range sbs.sorted {
		if !yield(bs) {
			return
		}
	}
}

// Add adds a new board to the sorted list if it is not a duplicate.
func (sbs *SortedBoardStates) Add(newBoard *BoardState) bool {
	insertIndex := len(sbs.sorted)
sortingLoop:
	for i, bs := range sbs.sorted {
		switch {
		case bs.Complexity() > newBoard.Complexity():
			continue
		case bs.Complexity() < newBoard.Complexity():
			insertIndex = i
			break sortingLoop
		}
		// do not add duplicates
		if boards.Equivalent(bs.board, newBoard.board) {
			return false
		}
	}
	sbs.sorted = slices.Insert(sbs.sorted, insertIndex, newBoard)
	return true
}

func (sbs *SortedBoardStates) AddAll(boards *SortedBoardStates) {
	// TODO: can optimize by sorting the input boards and deduping them
	for _, board := range boards.sorted {
		if !sbs.Add(board) {
			continue
		}
	}
}

func (sbs *SortedBoardStates) Results() []*solver.Result {
	return collections.MapSlice(sbs.sorted, (*BoardState).Result)
}
