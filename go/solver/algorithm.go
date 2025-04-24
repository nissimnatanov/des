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
	return []Algorithm{
		theOnlyChoice{},
		identifyPairs{},
		identifyTriplets{},
		newTrialAndError(),
	}

	/*
		missing for Solve that requires accurate leveling (but not solve fast nor prove):
			make_shared<SingleInSquare>(),
		    make_shared<SingleInRow>(),
		    make_shared<SingleInColumn>(),
	*/
}
