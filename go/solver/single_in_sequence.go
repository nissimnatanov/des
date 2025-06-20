package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
)

type singleInSequence struct {
}

func (a singleInSequence) String() string {
	return "Single in Sequence"
}

func (a singleInSequence) Run(_ context.Context, state AlgorithmState) Status {
	found := 0
	b := state.Board()
	for hints := b.Hints01(); hints != indexes.MinBitSet81; hints = b.Hints01() {
		foundNow := false
		for index := range b.Hints01().Indexes {
			allowed := b.AllowedValues(index).Values()
			switch len(allowed) {
			case 0:
				// found so far + 1 for No Solution detection
				state.AddStep(Step(a.String()), StepComplexityEasy, found+1)
				return StatusNoSolution
			case 1:
				if a.isSingleInAnySequence(b, index) {
					b.Set(index, allowed[0])
					found++
					foundNow = true
				}
			default:
				panic("Hint returned more than one allowed value")
			}
		}
		if !foundNow {
			break
		}
	}

	if found == 0 {
		return StatusUnknown
	}

	state.AddStep(Step(a.String()), StepComplexityEasy, found)
	return StatusSucceeded
}

func (a singleInSequence) isSingleInAnySequence(b *boards.Game, index int) bool {
	// Check if the missing value is the only allowed value in the row, column, or square
	missingOne := boards.SequenceSize - 1
	return b.RowValues(indexes.RowFromIndex(index)).Size() == missingOne ||
		b.ColumnValues(indexes.ColumnFromIndex(index)).Size() == missingOne ||
		b.SquareValues(indexes.SquareFromIndex(index)).Size() == missingOne
}
