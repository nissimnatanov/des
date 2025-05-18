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
	seqStatus := a.runSeqKind(ctx, state, indexes.RowSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, indexes.ColumnSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, indexes.SquareSequence)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	return status
}

func (a *theOnlyChoiceInSequence) runSeqKind(
	ctx context.Context, state AlgorithmState, seq func(seq int) indexes.Sequence,
) Status {
	for i := range boards.SequenceSize {
		if ctx.Err() != nil {
			return StatusError
		}
		status := a.runSeq(state, seq(i))
		if status != StatusUnknown {
			return status
		}
	}
	return StatusUnknown
}

func (a *theOnlyChoiceInSequence) runSeq(state AlgorithmState, seq indexes.Sequence) Status {
	b := state.Board()
	vs := values.Set(0)
	freeCells := a.indexWithAllowedCache[:0]
	for _, index := range seq {
		v := b.Get(index)
		if v == 0 {
			freeCells = append(freeCells, indexWithAllowed{
				index:   index,
				allowed: b.AllowedValues(index),
			})
		} else {
			vs = vs.With(v.AsSet())
		}
	}
	if vs.Size() == boards.SequenceSize {
		// all the values are set
		return StatusUnknown
	}

	// check if the missing values have a free cell in the sequence that allow them
	var found int
	for missingValue := range vs.Complement().Values {
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
