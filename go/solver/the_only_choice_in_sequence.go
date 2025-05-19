package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type indexWithAllowed struct {
	index   int
	allowed values.Set
}

type theOnlyChoiceInSequence struct {
	indexWithAllowedCache [indexes.BoardSequenceSize]indexWithAllowed
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
	freeCells := a.indexWithAllowedCache[:0]
	for index, allowed := range b.AllowedValuesIn(seq) {
		freeCells = append(freeCells, indexWithAllowed{
			index:   index,
			allowed: allowed,
		})
	}

	// check if the missing values have a free cell in the sequence that allow them
	var found int
missingValueLoop:
	for _, missingValue := range missingValues.Values() {
		foundIndex := -1
		for fi := range freeCells {
			if !freeCells[fi].allowed.Contains(missingValue) {
				continue
			}

			if foundIndex == -1 {
				// first host
				foundIndex = freeCells[fi].index
			} else {
				// second one, move to the next missing value
				continue missingValueLoop
			}
		}
		if foundIndex == -1 || !b.IsEmpty(foundIndex) {
			// either no host cell in sequence or same cell is forced to have two missing values,
			// fail the solution in both cases
			return StatusNoSolution
		}

		b.Set(foundIndex, missingValue)
		found++
	}

	if found > 0 {
		// we found at least one value to settle
		state.AddStep(Step(a.String()), StepComplexityMedium, found)
		return StatusSucceeded
	}

	return StatusUnknown
}
