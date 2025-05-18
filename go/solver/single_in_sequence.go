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

func (a singleInSequence) runSeqKind(
	ctx context.Context,
	state AlgorithmState,
	seq func(seq int) indexes.Sequence,
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

func (a singleInSequence) runSeq(state AlgorithmState, seq indexes.Sequence) Status {
	b := state.Board()
	vs := values.Set(0)
	freeCell := -1
	for _, index := range seq {
		v := b.Get(index)
		if v == 0 {
			// free cell
			if freeCell != -1 {
				return StatusUnknown
			}
			// first free cell
			freeCell = index
		} else {
			vs = vs.With(v.AsSet())
		}
	}
	if freeCell == -1 || vs.Size() != boards.SequenceSize-1 {
		return StatusUnknown
	}
	missingValue := vs.Complement().First()
	if !b.AllowedValues(freeCell).Contains(missingValue) {
		return StatusNoSolution
	}
	b.Set(freeCell, missingValue)
	state.AddStep(Step(a.String()), StepComplexityEasy, 1)
	return StatusSucceeded
}
