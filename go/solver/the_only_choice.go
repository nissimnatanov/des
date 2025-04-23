package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/board"
)

type theOnlyChoice struct {
}

func (a theOnlyChoice) Run(ctx context.Context, state AlgorithmState) Status {
	b := state.Board()
	status := StatusUnknown
	for i := range board.Size {
		if !b.IsEmpty(i) {
			continue
		}
		allowed := b.AllowedSet(i)
		switch allowed.Size() {
		case 0:
			return StatusNoSolution
		case 1:
			b.Set(i, allowed.At(0))
			state.AddStep(Step(a.String()), a.Complexity(), 1)
			status = StatusSucceeded
		}
	}
	return status
}

func (a theOnlyChoice) Complexity() StepComplexity {
	return StepComplexityMedium
}

func (a theOnlyChoice) String() string {
	return "The Only Choice"
}
