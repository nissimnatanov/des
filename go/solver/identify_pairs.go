package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type identifyPairs struct {
}

func (a identifyPairs) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	for index, allowed := range b.AllowedSets {
		if allowed.Size() != 2 {
			// this includes non-empty cells too (allowed set is empty)
			continue
		}
		peerIndex := a.findPeer(b, allowed, indexes.RelatedSequence(index))
		if peerIndex < 0 {
			continue
		}

		// found peer
		var eliminationCount int
		row := indexes.RowFromIndex(index)
		if row == indexes.RowFromIndex(peerIndex) {
			// same row
			if a.tryEliminate(b, index, peerIndex, allowed, indexes.RowSequence(row)) {
				eliminationCount++
			}
		}
		col := indexes.ColumnFromIndex(index)
		if col == indexes.ColumnFromIndex(peerIndex) {
			// same column
			if a.tryEliminate(b, index, peerIndex, allowed, indexes.ColumnSequence(col)) {
				eliminationCount++
			}
		}
		square := indexes.SquareFromIndex(index)
		if square == indexes.SquareFromIndex(peerIndex) {
			// same square
			if a.tryEliminate(b, index, peerIndex, allowed, indexes.SquareSequence(square)) {
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
	board *boards.Play, ignore1, ignore2 int,
	allowed values.Set, seq indexes.Sequence,
) bool {
	found := false
	for index := range seq.Indexes {
		if index == ignore1 || index == ignore2 || !board.IsEmpty(index) {
			continue
		}

		tempAllowed := board.AllowedSet(index)
		if values.Intersect(tempAllowed, allowed).Size() > 0 {
			// found a cell that we can remove values - turn them off
			board.DisallowSet(index, allowed)
			found = true
		}
	}

	return found
}

func (a identifyPairs) findPeer(board *boards.Play, allowed values.Set, seq indexes.Sequence) int {
	for peerIndex := range seq.Indexes {
		if !board.IsEmpty(peerIndex) {
			continue
		}
		peerAllowed := board.AllowedSet(peerIndex)
		if peerAllowed == allowed {
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
