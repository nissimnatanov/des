package solver

import (
	"context"
)

type theOnlyChoice struct {
}

// BenchmarkProve-10    	      40	  26921821 ns/op	 4437130 B/op	  219282 allocs/op
// BenchmarkProve-10    	      46	  25481993 ns/op	 3691628 B/op	  188216 allocs/op
// BenchmarkProve-10    	      48	  25315713 ns/op	 3374827 B/op	  168875 allocs/op
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
