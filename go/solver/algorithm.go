package solver

import (
	"context"

	"github.com/nissimnatanov/des/go/board"
)

type AlgorithmState interface {
	Board() board.Board

	Action() Action
	CurrentRecursionDepth() int
	MaxRecursionDepth() int
	AddStep(step Step, complexity StepComplexity, count int)
	MergeSteps(steps *StepStats)

	// recursiveRun is used to run the algorithm recursively
	recursiveRun(ctx context.Context, b board.Board) *Result
}

type Algorithm interface {
	String() string
	Run(ctx context.Context, state AlgorithmState) Status
}

// for now hardcoded algorithms, we can allow dynamic register for the algorithms later
var algorithms = []Algorithm{
	/*
			make_shared<SingleInSquare>(),
		    make_shared<SingleInRow>(),
		    make_shared<SingleInColumn>(),
	*/
	theOnlyChoice{},

	/*
	   // elimination algorithms
	   make_shared<IdentifyPairs>(),
	*/
	trialAndError{},
}
