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
	prover   *solver.Solver
}

type SolutionStateArgs struct {
	Solution *boards.Solution
	Rand     *random.Random
	Solver   *solver.Solver
	Prover   *solver.Solver
}

func NewSolutionState(args SolutionStateArgs) *SolutionState {
	if args.Rand == nil {
		args.Rand = random.New()
	}
	if args.Solution == nil {
		panic("solution must be provided")
	}
	if args.Solver == nil {
		args.Solver = solver.New(&solver.Options{Action: solver.ActionSolve})
	}
	if args.Prover == nil {
		args.Prover = solver.New(&solver.Options{Action: solver.ActionProve})
	}

	return &SolutionState{
		solver:   args.Solver,
		prover:   args.Prover,
		solution: args.Solution,
		rand:     args.Rand,
	}
}

func NewEnhanceBoardState(ctx context.Context, minLevel, maxLevel solver.Level, r *random.Random, board *boards.Game) (*SolutionState, *BoardState) {
	args := SolutionStateArgs{
		Rand:   r,
		Solver: solver.New(&solver.Options{Action: solver.ActionSolve}),
		Prover: solver.New(&solver.Options{Action: solver.ActionProve}),
	}
	res := args.Prover.Run(ctx, board)
	if res.Status != solver.StatusSucceeded {
		panic("cannot fine-tune unproven board")
	}
	args.Solution = res.Solutions.At(0)
	s := NewSolutionState(args)
	lr := LevelRange{
		Min: minLevel,
		Max: maxLevel,
	}
	lr = lr.WithDefaults()
	bs, _ := newBoardState(ctx, s, lr, board)
	return s, bs
}

func (s *SolutionState) InitialBoardState(ctx context.Context, levelRange LevelRange) *BoardState {
	return newSolutionBoardState(ctx, s, levelRange, s.solution)
}
