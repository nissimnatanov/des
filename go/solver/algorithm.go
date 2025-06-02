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
			// singleInSequence is the first since it is considered the easiest by
			// human to see the only missing value in a sequence (row/column/box)
			singleInSequence{},
			theOnlyAllowedValueInCell{},
			theOnlyChoiceInSequence{},
			identifyPairs{},
			identifyTriplets{},
			&squareToRowColumnConstraints{},
			&rowColToSquareConstraints{},
			newTrialAndError(),
		}
	case ActionSolveFast, ActionProve:
		return []Algorithm{
			// The only allowed value in cell is the most efficient algo, faster than
			// all the others. It is using the Hint01 method to identify cells
			// with 0 or 1 allowed value in O(1). It also covers what singleInSequence
			// would find, so no need to include this algo here.
			theOnlyAllowedValueInCell{},
			theOnlyChoiceInSequence{},
			&squareToRowColumnConstraints{},
			// recursion is faster than the identify pairs & triplets, and the row-col-2-square
			// algos - these algos are unfortunately too slow to beat the plain recursion
			// identifyPairs{},
			// identifyTriplets{},
			// &rowColToSquareConstraints{},
			newTrialAndError(),
		}
	default:
		panic("unknown action: " + action.String())
	}
}
