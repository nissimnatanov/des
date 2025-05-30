package internal

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/solver"
)

// State holds the initial state for the generator, including the level and
// base solution to use, then optional random number generator, solver, and prover.
type State struct {
	level    solver.Level
	solution *boards.Solution
	rand     *Random
	solver   *solver.Solver
	prover   *solver.Solver
}

type StateArgs struct {
	Level    solver.Level
	Solution *boards.Solution
	Rand     *Random
	Solver   *solver.Solver
	Prover   *solver.Solver
}

func NewState(args StateArgs) *State {
	if args.Rand == nil {
		args.Rand = NewRandom()
	}
	if args.Solution == nil {
		// to avoid circular dependencies, force caller to provide a solution
		// TODO: move solution generator to this internal pkg
		panic("internal GeneratorStateArgs must come with a solution")
	}
	if args.Level == solver.LevelUnknown {
		args.Level = solver.LevelEasy
	}
	if args.Solver == nil {
		args.Solver = solver.New(&solver.Options{Action: solver.ActionSolve})
	}
	if args.Prover == nil {
		args.Prover = solver.New(&solver.Options{Action: solver.ActionProve})
	}

	return &State{
		solver:   args.Solver,
		prover:   args.Prover,
		level:    args.Level,
		solution: args.Solution,
		rand:     args.Rand,
	}
}

func NewEnhanceBoardState(ctx context.Context, level solver.Level, r *Random, board *boards.Game) (*State, *BoardState) {
	args := StateArgs{
		Rand:   r,
		Level:  level,
		Solver: solver.New(&solver.Options{Action: solver.ActionSolve}),
		Prover: solver.New(&solver.Options{Action: solver.ActionProve}),
	}
	res := args.Prover.Run(ctx, board)
	if res.Status != solver.StatusSucceeded {
		panic("cannot fine-tune unproven board")
	}
	args.Solution = res.Solutions.At(0)
	s := NewState(args)
	bs := newBoardState(ctx, s, board)
	return s, bs
}

func (s *State) InitialBoardState(ctx context.Context) *BoardState {
	return newBoardState(ctx, s, s.solution)
}
