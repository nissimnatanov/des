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

const fastGenerationCap = solver.LevelVeryHard

type Options struct {
	Solution *boards.Solution
	MinLevel solver.Level
	// MaxLevel is optional, if not set it defaults to FromLevel.
	MaxLevel solver.Level
	// Count is the number of boards to generate per solution
	//
	// If not set, it defaults to 1 for fast-to-generate boards and an arbitrary number
	// for the slow ones.
	Count int
}

func (g *Generator) Generate(ctx context.Context, opts *Options) []*solver.Result {
	if opts == nil {
		opts = &Options{}
	} else {
		optsTmp := *opts
		opts = &optsTmp
	}

	if opts.Solution == nil {
		opts.Solution = GenerateSolution(g.r)
	}
	lr := internal.LevelRange{
		Min: opts.MinLevel,
		Max: opts.MaxLevel,
	}
	lr = lr.WithDefaults()

	state := internal.NewSolutionState(internal.SolutionStateArgs{
		Solution: opts.Solution,
		Rand:     g.r,
	})
	initState := state.InitialBoardState(ctx, internal.LevelRange{Min: opts.MinLevel, Max: opts.MaxLevel})
	if opts.MaxLevel > fastGenerationCap {
		return g.generateSlow(ctx, initState, opts.Count)
	}

	if opts.Count == 0 {
		opts.Count = 1 // default to 1 for fast generation
	}
	results := make([]*solver.Result, 0, opts.Count)
	for len(results) < opts.Count && ctx.Err() == nil {
		bs := g.generateFast(ctx, initState)
		results = append(results, bs.Result())
	}
	return results

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
func (g *Generator) generateFast(ctx context.Context, initState *internal.BoardState) *internal.BoardState {
	tries := 0
	start := time.Now()
	var stageStats GamePerStageStats
	for ctx.Err() == nil {
		tries++
		bs, stage := g.tryGenerateFastOnce(ctx, initState)
		stageStats.report(stage, bs != nil)
		if bs == nil {
			continue
		}
		elapsed := time.Since(start)
		Stats.reportGeneration(1, elapsed, int64(tries), stageStats)
		return bs
	}

	return nil
}

func (g *Generator) tryGenerateFastOnce(ctx context.Context, initState *internal.BoardState) (*internal.BoardState, int) {
	stage := 0
	bs := initState.Remove(ctx, internal.RemoveArgs{
		FreeCells:        45,
		BatchMinToRemove: 10,
		BatchMaxToRemove: 15,
		// first pass is usually needs to retry only once in hundreds of runs
		BatchMaxTries: 3,
	})
	if bs == nil {
		return nil, stage
	}

	if bs.Progress() == internal.InRangeStop {
		return bs, stage
	}
	stage++
	// remove the next bulk
	bs = bs.Remove(ctx, internal.RemoveArgs{
		FreeCells:        55,
		BatchMinToRemove: 2,
		BatchMaxToRemove: 4,
		BatchMaxTries:    40,
	})
	if bs == nil {
		return nil, stage
	}
	if bs.Progress() == internal.InRangeStop {
		return bs, stage
	}

	stage++
	// we have not reached the desired level yet, from this point remove one by one
	bs = bs.RemoveOneByOne(ctx)
	if bs == nil {
		return nil, stage
	}
	if bs.Progress() == internal.InRangeKeepGoing || bs.Progress() == internal.InRangeStop {
		return bs, stage
	}

	stage++
	if stage+1 != fastStageCount {
		panic("fast generation should end up with fastStageCount stages")
	}
	return nil, stage
}
