package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type squareToRowColumnConstraints struct {
}

func (a squareToRowColumnConstraints) Run(ctx context.Context, state AlgorithmState) Status {
	succeeded := false
	b := state.Board()
	freeBefore := b.FreeCellCount()
	for sqi := range boards.SequenceSize {
		sqVals := state.Board().SquareValues(sqi)
		if sqVals == values.FullSet {
			continue
		}
		sqSeq := indexes.SquareSequence(sqi)
		status := runConstraints(
			Step(a.String()), state, sqi, sqSeq,
			indexes.RowFromIndex, indexes.RowSequenceExcludeSquare)
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

		status = runConstraints(
			Step(a.String()), state, sqi, sqSeq,
			indexes.ColumnFromIndex, indexes.ColumnSequenceExcludeSquare)
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
	if succeeded {
		return StatusSucceeded
	}
	return StatusUnknown
}

func (a squareToRowColumnConstraints) Complexity() StepComplexity {
	return crossSequenceConstraintComplexity
}

func (a squareToRowColumnConstraints) String() string {
	return "Square to Row/Column Constraints"
}
