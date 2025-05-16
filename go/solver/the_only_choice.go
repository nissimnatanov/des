package solver

import (
	"context"
)

type theOnlyChoice struct {
}

func (a theOnlyChoice) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	found := 0
	for i, allowed := range b.AllowedSets {
		switch allowed.Size() {
		case 0:
			return StatusNoSolution
		case 1:
			b.Set(i, allowed.At(0))
			found++
		}
	}
	if found == 0 {
		return StatusUnknown
	}

	state.AddStep(Step(a.String()), a.Complexity(), found)
	return StatusSucceeded
}

func (a theOnlyChoice) Complexity() StepComplexity {
	return StepComplexityMedium
}

func (a theOnlyChoice) String() string {
	return "The Only Choice"
}
