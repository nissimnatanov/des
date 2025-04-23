package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/board"
	"github.com/nissimnatanov/des/go/board/indexes"
	"github.com/nissimnatanov/des/go/board/values"
)

type identifyPairs struct {
}

func (a identifyPairs) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	for index := range board.Size {
		if !b.IsEmpty(index) {
			continue
		}
		allowedValues := b.AllowedSet(index)
		if allowedValues.Size() != 2 {
			continue
		}
		peerIndex := a.findPeer(b, allowedValues, indexes.RelatedSequence(index))
		if peerIndex < 0 {
			continue
		}

		// found peer
		var eliminationCount int
		row := indexes.RowFromIndex(index)
		if row == indexes.RowFromIndex(peerIndex) {
			// same row
			if a.tryEliminate(b, index, peerIndex, allowedValues, indexes.RowSequence(row)) {
				eliminationCount++
			}
		}
		col := indexes.ColumnFromIndex(index)
		if col == indexes.ColumnFromIndex(peerIndex) {
			// same column
			if a.tryEliminate(b, index, peerIndex, allowedValues, indexes.ColumnSequence(col)) {
				eliminationCount++
			}
		}
		square := indexes.SquareFromIndex(index)
		if square == indexes.SquareFromIndex(peerIndex) {
			// same square
			if a.tryEliminate(b, index, peerIndex, allowedValues, indexes.SquareSequence(square)) {
				eliminationCount++
			}
		}
		if eliminationCount > 0 {
			state.AddStep(Step(a.String()), StepComplexityHard, eliminationCount)
			return StatusSucceeded
		}
	}
	return StatusUnknown
}

func (a identifyPairs) tryEliminate(
	board board.Board, ignore1, ignore2 int,
	allowedValues values.Set, indexes indexes.Sequence,
) bool {
	found := false
	for temp := 0; temp < 9; temp++ {
		index := indexes.Get(temp)
		if index == ignore1 || index == ignore2 || !board.IsEmpty(index) {
			continue
		}

		tempAllowedValues := board.AllowedSet(index)
		if values.Intersect(tempAllowedValues, allowedValues).Size() > 0 {
			// found a cell that we can remove values - turn them off
			board.DisallowSet(index, allowedValues)
			found = true
		}
	}

	return found
}

func (a identifyPairs) findPeer(board board.Board, allowedValues values.Set, indexes indexes.Sequence) int {
	for pi := range indexes.Size() {
		peerIndex := indexes.Get(pi)
		if !board.IsEmpty(peerIndex) {
			continue
		}
		peerAllowedValues := board.AllowedSet(peerIndex)
		if peerAllowedValues == allowedValues {
			return peerIndex
		}
	}
	return -1
}

func (a identifyPairs) Complexity() StepComplexity {
	return StepComplexityHard
}

func (a identifyPairs) String() string {
	return "Identify Pairs"
}
