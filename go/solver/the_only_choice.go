package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/board"
)

type theOnlyChoice struct {
}

func (a theOnlyChoice) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	for i := range board.Size {
		if !b.IsEmpty(i) {
			continue
		}
		switch b.AllowedSet(i).Size() {
		case 0:
			return StatusNoSolution
		case 1:
			b.Set(i, b.AllowedSet(i).At(0))
			state.AddStep(Step(a.String()), a.Complexity(), 1)
			return StatusSucceeded
		}
	}
	return StatusUnknown
}

func (a theOnlyChoice) Complexity() StepComplexity {
	return StepComplexityMedium
}

func (a theOnlyChoice) String() string {
	return "The Only Choice"
}
