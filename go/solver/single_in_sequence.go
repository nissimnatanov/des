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
	seqStatus := a.runSeqKind(state, indexes.RowSequence, (*boards.Game).RowValues)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(state, indexes.ColumnSequence, (*boards.Game).ColumnValues)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	seqStatus = a.runSeqKind(state, indexes.SquareSequence, (*boards.Game).SquareValues)
	if seqStatus != StatusUnknown {
		return seqStatus
	}
	return status
}

func (a singleInSequence) runSeqKind(
	state AlgorithmState,
	seq func(seq int) indexes.Sequence,
	seqValues func(b *boards.Game, seq int) values.Set,
) Status {
	for si := range boards.SequenceSize {
		seqValues := seqValues(state.Board(), si)
		if seqValues.Size() != (boards.SequenceSize - 1) {
			continue
		}
		missingValue := seqValues.Complement().First()
		status := a.setMissingValue(state, seq(si), missingValue)
		if status != StatusUnknown {
			return status
		}
	}
	return StatusUnknown
}

func (a singleInSequence) setMissingValue(
	state AlgorithmState,
	seq indexes.Sequence,
	missingValue values.Value) Status {
	b := state.Board()
	for index, allowed := range b.AllowedValuesIn(seq) {
		if !allowed.Contains(missingValue) {
			return StatusNoSolution
		}
		b.Set(index, missingValue)
		state.AddStep(Step(a.String()), StepComplexityEasy, 1)
		return StatusSucceeded
	}
	return StatusError
}
