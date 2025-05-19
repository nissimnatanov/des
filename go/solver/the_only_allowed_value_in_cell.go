package solver

import (
	"context"
)

type theOnlyAllowedValueInCell struct {
}

func (a theOnlyAllowedValueInCell) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	found := 0
	for index := b.Hint01(); index >= 0; index = b.Hint01() {
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
