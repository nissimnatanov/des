package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type singleInSequence struct {
}

func (a singleInSequence) String() string {
	return "Single in Sequence"
}

func (a singleInSequence) Run(ctx context.Context, state AlgorithmState) Status {
	freeCellsCache := [boards.SequenceSize]indexWithAllowed{}
	status := StatusUnknown
	seqStatus := a.runSeqKind(ctx, state, indexes.RowSequence, freeCellsCache[:])
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, indexes.ColumnSequence, freeCellsCache[:])
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(ctx, state, indexes.SquareSequence, freeCellsCache[:])
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	return status
}

func (a singleInSequence) runSeqKind(
	ctx context.Context, state AlgorithmState, seq func(seq int) indexes.Sequence,
	freeCellsCache []indexWithAllowed,
) Status {
	for i := range boards.SequenceSize {
		if ctx.Err() != nil {
			return StatusError
		}
		status := a.runSeq(state, seq(i), freeCellsCache)
		if status != StatusUnknown {
			return status
		}
	}
	return StatusUnknown
}

func (a singleInSequence) runSeq(
	state AlgorithmState, seq indexes.Sequence, freeCellsCache []indexWithAllowed,
) Status {
	b := state.Board()
	vs := values.Set(0)
	freeCells := freeCellsCache[:0]
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
