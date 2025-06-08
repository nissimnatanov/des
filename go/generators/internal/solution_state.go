package internal

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/solver"
)

// SolutionState holds the initial state for the generator, including the level and
// base solution to use, then optional random number generator, solver, and prover.
type SolutionState struct {
	solution *boards.Solution
	rand     *random.Random
	solver   *solver.Solver
}

type SolutionStateArgs struct {
	Solution *boards.Solution
	Rand     *random.Random
	Solver   *solver.Solver
}

func NewSolutionState(args SolutionStateArgs) *SolutionState {
	if args.Rand == nil {
		args.Rand = random.New()
	}
	if args.Solution == nil {
		panic("solution must be provided")
	}
	if args.Solver == nil {
		args.Solver = solver.New()
	}

	return &SolutionState{
		solver:   args.Solver,
		solution: args.Solution,
		rand:     args.Rand,
	}
}

func (s *SolutionState) InitialBoardState(ctx context.Context, levelRange LevelRange) *BoardState {
	return newSolutionBoardState(ctx, s, levelRange, s.solution)
}
