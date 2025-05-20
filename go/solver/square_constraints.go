package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type squareToRowColumnConstraints struct {
	valRowCache   [10]int
	valueColCache [10]int
}

func (a *squareToRowColumnConstraints) Run(ctx context.Context, state AlgorithmState) Status {
	succeeded := false
	b := state.Board()
	freeBefore := b.FreeCellCount()
	for sqi := range boards.SequenceSize {
		sqVals := state.Board().SquareValues(sqi)
		if sqVals == values.FullSet {
			continue
		}
		status := a.runSquareConstraints(state, sqi, indexes.SquareSequence(sqi))
		switch status {
		case StatusUnknown:
			// keep looking
		case StatusSucceeded:
			// this is elimination only algorithm, it is more efficient to keep looking
			// for all other options than bailing out now
			succeeded = true
			if b.FreeCellCount() < freeBefore || state.Action().LevelRequested() {
				// If a value was set, or if an accurate level is requested, stop now to
				// allow lower-level algorithms to do the job.
				// If, however, we only eliminated a value and level is not needed, it is
				// more efficient to keep looking for other options within the same algo.
				return StatusSucceeded
			}
		default:
			// if we eliminated the only candidate in the other square, bail out
			return status
		}
	}
	if succeeded {
		return StatusSucceeded
	}
	return StatusUnknown
}

func (a *squareToRowColumnConstraints) runSquareConstraints(state AlgorithmState, sqi int, sqSeq indexes.Sequence) Status {
	b := state.Board()
	valRows := a.valRowCache[:10]
	valCols := a.valueColCache[:10]
	for i := range 10 {
		valRows[i] = -1
		valCols[i] = -1
	}

	foundCandidates := 0
	for index, allowed := range b.AllowedValuesIn(sqSeq) {
		row := indexes.RowFromIndex(index)
		col := indexes.ColumnFromIndex(index)

		for _, v := range allowed.Values() {
			if valRows[v] == -1 {
				valRows[v] = row
				foundCandidates++
			} else if valRows[v] >= 0 && valRows[v] != row {
				valRows[v] = -2
				foundCandidates--
			}
			if valCols[v] == -1 {
				valCols[v] = col
				foundCandidates++
			} else if valCols[v] >= 0 && valCols[v] != col {
				valCols[v] = -2
				foundCandidates--
			}
		}
	}

	// let's see if any value must be present in a specific row or column
	if foundCandidates == 0 {
		return StatusUnknown
	}
	// now, let's try to eliminate the same value in all other squares
	eliminateCount := 0
	for v := values.Value(1); v <= 9; v++ {
		if valRows[v] >= 0 {
			rowSeqNotSharedWithSquare := indexes.RowSequenceExcludeSquare(valRows[v], sqi)
			status := eliminateInSequence(b, v, rowSeqNotSharedWithSquare)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				eliminateCount++
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
		}
		if valCols[v] >= 0 {
			colSeqNotSharedWithSquare := indexes.ColumnSequenceExcludeSquare(valCols[v], sqi)
			status := eliminateInSequence(b, v, colSeqNotSharedWithSquare)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				eliminateCount++
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
		}
	}
	if eliminateCount > 0 {
		state.AddStep(Step(a.String()), a.Complexity(), eliminateCount)
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
