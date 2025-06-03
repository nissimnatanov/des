package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type rowColToSquareConstraints struct {
}

func (a rowColToSquareConstraints) Run(ctx context.Context, state AlgorithmState) Status {
	succeeded := false
	b := state.Board()
	freeBefore := b.FreeCellCount()
	for rci := range boards.SequenceSize {
		rowVals := b.RowValues(rci)
		if rowVals != values.FullSet {
			status := runConstraints(
				Step(a.String()), state, rci, indexes.RowSequence(rci),
				indexes.SquareFromIndex, indexes.SquareSequenceExcludeRow)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				succeeded = true
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
			if succeeded && (b.FreeCellCount() < freeBefore || state.Action().LevelRequested()) {
				// If a value was set, or if an accurate level is requested, stop now to
				// allow lower-level algorithms to do the job.
				// If, however, we only eliminated a value and level is not needed, it is
				// more efficient to keep looking for other options within the same algo.
				return StatusSucceeded
			}
		}
		colVals := b.ColumnValues(rci)
		if colVals != values.FullSet {
			status := runConstraints(
				Step(a.String()), state, rci, indexes.ColumnSequence(rci),
				indexes.SquareFromIndex, indexes.SquareSequenceExcludeColumn)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				succeeded = true
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
			if succeeded && (b.FreeCellCount() < freeBefore || state.Action().LevelRequested()) {
				// If a value was set, or if an accurate level is requested, stop now to
				// allow lower-level algorithms to do the job.
				// If, however, we only eliminated a value and level is not needed, it is
				// more efficient to keep looking for other options within the same algo.
				return StatusSucceeded
			}
		}
	}
	if succeeded {
		return StatusSucceeded
	}
	return StatusUnknown
}

func (a rowColToSquareConstraints) Complexity() StepComplexity {
	return crossSequenceConstraintComplexity
}

func (a rowColToSquareConstraints) String() string {
	return "Row/Column to Square Constraints"
}
