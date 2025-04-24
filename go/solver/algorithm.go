package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
)

type AlgorithmState interface {
	Board() *boards.Game

	Action() Action
	CurrentRecursionDepth() int
	MaxRecursionDepth() int
	AddStep(step Step, complexity StepComplexity, count int)
	MergeSteps(steps *StepStats)

	// recursiveRun is used to run the algorithm recursively
	recursiveRun(ctx context.Context, b *boards.Game) *Result
}

type Algorithm interface {
	String() string
	Run(ctx context.Context, state AlgorithmState) Status
}

// for now hardcoded algorithms, we can allow dynamic register for the algorithms later
func GetAlgorithms(action Action) []Algorithm {
	// for now - same algos for all actions
	switch action {
	case ActionSolve:
		return []Algorithm{
			// singleInSequence is first since it is considered the easiest by human to see
			singleInSequence{},
			theOnlyChoice{},
			identifyPairs{},
			identifyTriplets{},
			newTrialAndError(),
		}
	case ActionSolveFast, ActionProve:
		return []Algorithm{
			theOnlyChoice{},
			identifyPairs{},
			// singleInSequence is slower that the first two, but faster than the triplet and recursion
			singleInSequence{},
			identifyTriplets{},
			newTrialAndError(),
		}
	default:
		panic("unknown action: " + action.String())
	}
}
