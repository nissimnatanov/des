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
	seqStatus := a.runSeqKind(ctx, state, b.RowValues, indexes.RowSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, b.ColumnValues, indexes.ColumnSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, b.SquareValues, indexes.SquareSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	return status
}

func (a *theOnlyChoiceInSequence) runSeqKind(
	ctx context.Context,
	state AlgorithmState,
	seqValues func(seq int) values.Set,
	seq func(seq int) indexes.Sequence,
) Status {
	for si := range boards.SequenceSize {
		if ctx.Err() != nil {
			return StatusError
		}
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
	for missingValue := range missingValues.Values {
		freeIndex := -1
		freeIndex2 := -1
		for _, freeCell := range freeCells {
			if !freeCell.allowed.Contains(missingValue) {
				continue
			}

			if freeIndex == -1 {
				// first free cell
				freeIndex = freeCell.index
			} else if freeIndex2 == -1 {
				// second one, we can stop now
				freeIndex2 = freeCell.index
				break
			}
		}
		switch {
		case freeIndex == -1:
			// no host cell in sequence - fail the solution
			return StatusNoSolution
		case freeIndex2 != -1:
			// two cells to allow the value, try next one
			continue
		}
		// we have exactly one cell that allows the value
		if !b.IsEmpty(freeIndex) {
			// same cell is forced to have two missing values
			return StatusNoSolution
		}

		b.Set(freeIndex, missingValue)
		found++
	}

	if found > 0 {
		// we found at least one value to settle
		state.AddStep(Step(a.String()), StepComplexityEasy, found)
		return StatusSucceeded
	}

	return StatusUnknown
}
