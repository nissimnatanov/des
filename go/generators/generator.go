package generators

import (
	"context"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/generators/solution"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/solver"
)

func New(opts *Options) *Generator {
	if opts == nil {
		opts = &Options{}
	}
	var r *random.Random
	if opts.RandSeed == 0 {
		r = random.New()
	} else {
		r = random.WithSeed(opts.RandSeed)
	}
	lr := internal.LevelRange{
		Min: opts.MinLevel,
		Max: opts.MaxLevel,
	}
	lr = lr.WithDefaults()
	count := opts.Count
	if count < 0 {
		count = 0 // all available boards, at least one
	}
	solProvider := opts.SolutionProvider
	if solProvider == nil {
		solProvider = func() *boards.Solution {
			return solution.Generate(r)
		}
	}

	g := &Generator{
		r:           r,
		lr:          lr,
		count:       count,
		solProvider: solProvider,
	}
	return g
}

type Generator struct {
	r     *random.Random
	lr    internal.LevelRange
	count int

	solProvider func() *boards.Solution
	// if set, it will be used to enhance the board
	// enhanceBoard *boards.Game
}

func (g *Generator) Seed() int64 {
	return g.r.Seed()
}

type Options struct {
	RandSeed int64 // optional, if 0, a new random seed will be generated
	MinLevel solver.Level
	// MaxLevel is optional, if not set it defaults to FromLevel.
	MaxLevel solver.Level
	// Count is the number of boards to generate per solution
	//
	// If not set, it defaults to 1 for fast-to-generate boards and an arbitrary number
	// for the slow ones.
	Count int

	// optional
	SolutionProvider func() *boards.Solution
}

func (g *Generator) Generate(ctx context.Context) []*solver.Result {
	state := internal.NewSolutionState(internal.SolutionStateArgs{
		Solution: g.solProvider(),
		Rand:     g.r,
	})
	initState := state.InitialBoardState(ctx, g.lr)
	if g.lr.Max > internal.FastGenerationCap {
		return g.generateSlow(ctx, initState, g.count)
	}

	count := g.count
	if count == 0 {
		count = 1 // default to 1 for fast generation
	}
	results := make([]*solver.Result, 0, count)
	for len(results) < count && ctx.Err() == nil {
		bs := g.generateFast(ctx, initState)
		results = append(results, bs.Result())
	}
	return results
}

/*
// Enhance tries removing values from the existing board until it reaches the desired level.
func (g *Generator) Enhance(ctx context.Context, board *boards.Game, level solver.Level) *solver.Result {
	bs := internal.NewEnhanceBoardState(ctx, level, g.r, board)
	if level <= internal.FastGenerationCap {
		return g.generateFast(ctx, bs)
	}

	return g.generateSlow(ctx, bs)
}*/

// generateFast for lower levels
func (g *Generator) generateFast(ctx context.Context, initState *internal.BoardState) *internal.BoardState {
	tries := 0
	start := time.Now()
	var stageStats internal.GamePerStageStats
	for ctx.Err() == nil {
		tries++
		bs, stage := g.tryGenerateFastOnce(ctx, initState)
		stageStats.Report(stage, bs != nil)
		if bs == nil {
			continue
		}
		elapsed := time.Since(start)
		internal.Stats.ReportGeneration(1, elapsed, int64(tries), stageStats)
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
	if stage+1 != internal.FastStageCount {
		panic("fast generation should end up with fastStageCount stages")
	}
	return nil, stage
}
