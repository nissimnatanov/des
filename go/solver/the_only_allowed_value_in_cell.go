package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards/indexes"
)

type theOnlyAllowedValueInCell struct {
}

func (a theOnlyAllowedValueInCell) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	if state.Action().LevelRequested() {
		index := b.Hint01()
		if index < 0 {
			return StatusUnknown
		}
		allowed := b.AllowedValues(index).Values()
		switch len(allowed) {
		case 0:
			return StatusNoSolution
		case 1:
			b.Set(index, allowed[0])
			state.AddStep(Step(a.String()), a.Complexity(), 1)
			return StatusSucceeded
		}

		panic("Hint returned more than one allowed value")
	}

	found := 0
	for hints := b.Hints01(); hints != indexes.MinBitSet81; hints = b.Hints01() {
		for index := range b.Hints01().Indexes {
			allowed := b.AllowedValues(index).Values()
			switch len(allowed) {
			case 0:
				return StatusNoSolution
			case 1:
				b.Set(index, allowed[0])
				found++
			default:
				panic("Hint returned more than one allowed value")
			}
		}
	}
	if found == 0 {
		return StatusUnknown
	}

	state.AddStep(Step(a.String()), a.Complexity(), found)
	return StatusSucceeded
}

func (a theOnlyAllowedValueInCell) Complexity() StepComplexity {
	return StepComplexityMedium
}

func (a theOnlyAllowedValueInCell) String() string {
	return "The Only Allowed Value in Cell"
}
