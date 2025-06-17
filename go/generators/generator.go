package generators

import (
	"context"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/values"
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
	count := max(opts.Count, 0)
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
		onNewResult: opts.OnNewResult,
		solver:      solver.New(),
	}
	return g
}

// Generator is single-threaded!
type Generator struct {
	r     *random.Random
	lr    internal.LevelRange
	count int

	solProvider func() *boards.Solution
	onNewResult func(*solver.Result)

	// cache the solver instances to preserve its caches and keep the memory
	// profiling cleaner
	solver *solver.Solver

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

	// OnNewResult is an optional callback that will be called for each new board's result generated.
	// Note: it can be called multiple times for the same board and it can be called more than
	// the requested Count time (deduplication happens at the end of the generation process).
	OnNewResult func(*solver.Result)
}

func (g *Generator) newInitialBoardState(ctx context.Context, withCache bool) *internal.BoardState {
	state := internal.NewSolutionState(internal.SolutionStateArgs{
		Solution:  g.solProvider(),
		Rand:      g.r,
		Solver:    g.solver,
		WithCache: withCache,
	})
	return state.InitialBoardState(ctx, g.lr)
}

func (g *Generator) Generate(ctx context.Context) []*solver.Result {
	if g.lr.Max > internal.FastGenerationCap {
		return g.generateSlow(ctx)
	}

	// cache degrades generation performance for fast boards
	initState := g.newInitialBoardState(ctx, false)
	count := g.count
	if count == 0 {
		count = 1 // default to 1 for fast generation
	}
	results := make([]*solver.Result, 0, count)
	for len(results) < count && ctx.Err() == nil {
		bs := g.removeSingleValue(ctx, initState)
		bs = g.generateFast(ctx, bs)
		results = append(results, bs.Result())
		if g.onNewResult != nil {
			g.onNewResult(bs.Result())
		}
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
		if bs == nil {
			stageStats.Report(0, stage)
			continue
		}

		stageStats.Report(1, stage)
		elapsed := time.Since(start)
		internal.Stats.ReportGeneration(1, elapsed, int64(tries), stageStats, initState.SolutionState().Cache().Stats())
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
	bs = bs.RemoveOneByOne(ctx, solver.MaxFreeCellsForValidBoard)
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

// removeSingleValue must be called on full board only
func (g *Generator) removeSingleValue(ctx context.Context, initState *internal.BoardState) *internal.BoardState {
	minLevel := initState.DesiredLevelRange().Min
	v := values.Value(g.r.Intn(9) + 1)
	var removeCount int
	switch {
	case minLevel <= solver.LevelMedium:
		// for perf reasons mostly, not much benefit otherwise
		removeCount = 5
	case minLevel < solver.LevelVeryHard:
		switch g.r.Intn(3) {
		case 0:
			removeCount = 5
		case 1:
			removeCount = 6
		default:
			removeCount = 7
		}
	case minLevel < solver.LevelNightmare:
		switch g.r.Intn(6) {
		case 0:
			removeCount = 6
		case 1, 2:
			removeCount = 7
		default:
			removeCount = 8
		}
	default:
		switch g.r.Intn(6) {
		case 0:
			removeCount = 7
		case 1, 2:
			removeCount = 8
		default:
			removeCount = 9
		}
	}

	return initState.RemoveVal(ctx, v, removeCount)
}
