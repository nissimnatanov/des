package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type theOnlyChoiceInSequence struct {
}

func (a *theOnlyChoiceInSequence) String() string {
	return "The Only Choice in Sequence"
}

func (a *theOnlyChoiceInSequence) Run(ctx context.Context, state AlgorithmState) Status {
	status := StatusUnknown
	b := state.Board()
	seqStatus := a.runSeqKind(state, b.RowValues, indexes.RowSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(state, b.ColumnValues, indexes.ColumnSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(state, b.SquareValues, indexes.SquareSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	return status
}

func (a *theOnlyChoiceInSequence) runSeqKind(
	state AlgorithmState,
	seqValues func(seq int) values.Set,
	seq func(seq int) indexes.Sequence,
) Status {
	for si := range boards.SequenceSize {
		vs := seqValues(si)
		if vs == values.FullSet {
			continue
		}
		status := a.runSeq(state, vs.Complement(), seq(si))
		if status != StatusUnknown {
			return status
		}
	}
	return StatusUnknown
}

func (a *theOnlyChoiceInSequence) runSeq(
	state AlgorithmState,
	missingValues values.Set,
	seq indexes.Sequence,
) Status {
	b := state.Board()
	var uniqueValues values.Set
	for _, allowed := range b.AllowedValuesIn(seq) {
		newUniqueValues := values.Intersect(missingValues, allowed)
		missingValues = missingValues.Without(allowed)

		// values that were unique but appear in the current cell are no
		// longer unique
		uniqueValues = uniqueValues.Without(allowed)
		// if this cell covers new values from the remained missing values,
		// let's add them back to the unique values
		uniqueValues = uniqueValues.With(newUniqueValues)
	}

	if !missingValues.IsEmpty() {
		// we have values that are not covered anywhere
		return StatusNoSolution
	}
	if uniqueValues.IsEmpty() {
		// we do not have any unique values in the sequence
		return StatusUnknown
	}

	// we have at least one unique value in the sequence, lets find its cell
	for index, allowed := range b.AllowedValuesIn(seq) {
		uniqueInCell := values.Intersect(uniqueValues, allowed)
		switch uniqueInCell.Size() {
		case 0:
			// wrong cell, keep going
			continue
		case 1:
			b.Set(index, uniqueInCell.First())
		default:
			// we cannot put two unique values in the same cell
			return StatusNoSolution
		}
	}

	state.AddStep(Step(a.String()), StepComplexityMedium, uniqueValues.Size())
	return StatusSucceeded
}
