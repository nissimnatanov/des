package generators

import (
	"context"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/solver"
)

func New() *Generator {
	return &Generator{
		r: internal.NewRandom(),
	}
}

type Generator struct {
	r *internal.Random
}

// TODO: adjust the bar
const fastGenerationCap = solver.LevelBlackHole

func (g *Generator) Generate(ctx context.Context, level solver.Level, sol *boards.Solution) *solver.Result {
	if sol == nil {
		sol = GenerateSolution(g.r)
	}
	state := internal.NewState(internal.StateArgs{
		Level:    level,
		Solution: sol,
		Rand:     g.r,
	})
	if level <= fastGenerationCap {
		return g.generateFast(ctx, state)
	}

	return g.generateSlow(ctx, state)
}

/*
// Enhance tries removing values from the existing board until it reaches the desired level.
func (g *Generator) Enhance(ctx context.Context, board *boards.Game, level solver.Level) *solver.Result {
	bs := internal.NewEnhanceBoardState(ctx, level, g.r, board)
	if level <= fastGenerationCap {
		return g.generateFast(ctx, bs)
	}

	return g.generateSlow(ctx, bs)
}*/

// generate (for lower levels)
func (g *Generator) generateFast(ctx context.Context, state *internal.State) *solver.Result {
	start := state.InitialBoardState(ctx)
	for ctx.Err() == nil {
		bs := start.Remove(ctx, internal.RemoveArgs{
			FreeAtLeast: 32,
			MinToRemove: 3,
			MaxToRemove: 8,
			MaxRetries:  15,
		})
		if bs == nil {
			continue
		}
		if bs.Progress() == internal.AtLevelStop {
			return bs.Result()
		}

		// remove the next bulk
		bs = bs.Remove(ctx, internal.RemoveArgs{
			FreeAtLeast: 48,
			MinToRemove: 1,
			MaxToRemove: 3,
			MaxRetries:  25,
		})
		if bs == nil {
			continue
		}
		if bs.Progress() == internal.AtLevelStop {
			return bs.Result()
		}

		// we have not reached the desired level yet, from this point remove one by one
		bs = bs.RemoveOneByOne(ctx)
		if bs == nil {
			continue
		}
		if bs.Progress() == internal.AtLevelKeepGoing || bs.Progress() == internal.AtLevelStop {
			return bs.Result()
		}
	}

	return nil
}
