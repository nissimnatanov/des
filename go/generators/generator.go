package generators

import (
	"context"
	"time"

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

const fastGenerationCap = solver.LevelEvil

type Options struct {
	Solution *boards.Solution
}

func (g *Generator) Generate(ctx context.Context, level solver.Level, opts *Options) *solver.Result {
	var sol *boards.Solution
	if opts != nil {
		sol = opts.Solution
	}
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

// generateFast for lower levels
func (g *Generator) generateFast(ctx context.Context, state *internal.State) *solver.Result {
	initState := state.InitialBoardState(ctx)
	tries := 0
	start := time.Now()
	var stages [4]GameStageStats

	for ctx.Err() == nil {
		tries++
		res, stage := g.tryGenerateFastOnce(ctx, initState)
		for s := range stage + 1 {
			stages[s].Total++
		}
		if res == nil {
			stages[stage].Failed++
			continue
		}
		stages[stage].Succeeded++
		elapsed := time.Since(start)
		Stats.reportOneGeneration(elapsed, int64(tries), res.Steps.Complexity, stages[:])
		return res
	}

	return nil
}

func (g *Generator) tryGenerateFastOnce(ctx context.Context, initState *internal.BoardState) (*solver.Result, int) {
	stage := 0
	bs := initState.Remove(ctx, internal.RemoveArgs{
		FreeAtLeast:      45,
		BatchMinToRemove: 10,
		BatchMaxToRemove: 15,
		// first pass is usually needs to retry only once in hundreds of runs
		BatchMaxTries: 3,
	})
	if bs == nil {
		return nil, stage
	}

	if bs.Progress() == internal.AtLevelStop {
		return bs.Result(), stage
	}
	stage++
	// remove the next bulk
	bs = bs.Remove(ctx, internal.RemoveArgs{
		FreeAtLeast:      55,
		BatchMinToRemove: 3,
		BatchMaxToRemove: 5,
		BatchMaxTries:    30,
	})
	if bs == nil {
		return nil, stage
	}
	if bs.Progress() == internal.AtLevelStop {
		return bs.Result(), stage
	}

	stage++
	// we have not reached the desired level yet, from this point remove one by one
	bs = bs.RemoveOneByOne(ctx)
	if bs == nil {
		return nil, stage
	}
	if bs.Progress() == internal.AtLevelKeepGoing || bs.Progress() == internal.AtLevelStop {
		return bs.Result(), stage
	}

	stage++
	return nil, stage
}
