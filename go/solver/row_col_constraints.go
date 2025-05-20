package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type rowColToSquareConstraints struct {
	valSqCache [10]int
}

func (a *rowColToSquareConstraints) Run(ctx context.Context, state AlgorithmState) Status {
	succeeded := false
	b := state.Board()
	freeBefore := b.FreeCellCount()
	for rci := range boards.SequenceSize {
		rowVals := b.RowValues(rci)
		if rowVals != values.FullSet {
			status := a.runRowColConstraints(
				state, rci, indexes.RowSequence(rci), indexes.SquareSequenceExcludeRow)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				succeeded = true
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
		}
		colVals := b.ColumnValues(rci)
		if colVals != values.FullSet {
			status := a.runRowColConstraints(
				state, rci, indexes.ColumnSequence(rci), indexes.SquareSequenceExcludeColumn)
			switch status {
			case StatusUnknown:
				// keep looking
			case StatusSucceeded:
				succeeded = true
			default:
				// if we eliminated the only candidate in the other square, bail out
				return status
			}
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

func (a *rowColToSquareConstraints) runRowColConstraints(
	state AlgorithmState,
	// either row or column
	rci int,
	// the sequence of board indexes in this row/column
	rcSeq indexes.Sequence,
	sqSeqNotSharedWithRC func(sq int, rc int) indexes.Sequence,
) Status {
	b := state.Board()
	valSq := a.valSqCache[:10]
	for i := range 10 {
		valSq[i] = -1
	}

	foundCandidates := 0
	for index, allowed := range b.AllowedValuesIn(rcSeq) {
		sq := indexes.SquareFromIndex(index)
		for _, v := range allowed.Values() {
			if valSq[v] == -1 {
				valSq[v] = sq
				foundCandidates++
			} else if valSq[v] >= 0 && valSq[v] != sq {
				valSq[v] = -2
				foundCandidates--
			}
		}
	}

	// let's see if any value must be present in a specific row or column
	if foundCandidates == 0 {
		return StatusUnknown
	}
	// now, let's try to eliminate the same value in all other rows/columns
	eliminateCount := 0
	for v := values.Value(1); v <= 9; v++ {
		if valSq[v] >= 0 {
			sqSeqNotSharedWithRow := sqSeqNotSharedWithRC(valSq[v], rci)
			status := eliminateInSequence(b, v, sqSeqNotSharedWithRow)
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

func (a rowColToSquareConstraints) Complexity() StepComplexity {
	return crossSequenceConstraintComplexity
}

func (a rowColToSquareConstraints) String() string {
	return "Row/Column to Square Constraints"
}
