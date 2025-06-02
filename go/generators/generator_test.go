package generators_test

import (
	"runtime/debug"
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestGeneratorFast(t *testing.T) {
	const loop = 10
	ctx := t.Context()
	for level := solver.LevelEasy; level <= internal.FastGenerationCap; level++ {
		t.Run(level.String(), func(t *testing.T) {
			for range loop {
				g := generators.New(&generators.Options{MinLevel: level, MaxLevel: level})
				rs := g.Generate(ctx)
				assert.Check(t, cmp.Len(rs, 1), "expected exactly one result for level %s, got %d", level, len(rs))
				for _, res := range rs {
					assert.Assert(t, res != nil)
					assert.Equal(t, res.Status, solver.StatusSucceeded)
					assert.Check(t, res.Steps.Level == level, "expected level %s, got %s", level, res.Steps.Level)
				}
			}
		})
	}
}

func TestGeneratorRangeWithSlowMax(t *testing.T) {
	const loop = 10
	ctx := t.Context()
	for range loop {
		min := solver.LevelHard
		max := solver.LevelDarkEvil // evil requires slow generation
		g := generators.New(&generators.Options{MinLevel: min, MaxLevel: max, Count: 3})
		rs := g.Generate(ctx)
		assert.Check(t, cmp.Len(rs, 3), "expected three results, got %d", len(rs))
		for _, res := range rs {
			assert.Assert(t, res != nil)
			assert.Equal(t, res.Status, solver.StatusSucceeded)
			assert.Check(t, res.Steps.Level >= min && res.Steps.Level <= max,
				"expected level between %s and %s, got %s", min, max, res.Steps.Level)
		}
	}
}

// Initial state (bulks of 100):
// * 18	  57867785 ns/op	 7639343 B/op	   50437 allocs/op
// * Generations: 1800, ~Elapsed: 577.917µs, ~Retries: 1.000,
//   Stages: [{1800 452 0} {1348 1048 0} {300 300 0}]

func BenchmarkEasyOrMedium(b *testing.B) {
	runBenchmark(b, solver.LevelEasy, solver.LevelMedium, 100)
}

// Initial state (bulks of 10):
// * 22	  59968347 ns/op	 6069760 B/op	   41416 allocs/op
// * Generations: 220, ~Elapsed: 5.992518ms, ~Retries: 2.205,
//   Stages: [{220 0 0} {220 42 0} {178 178 6}]

func BenchmarkHardOrVeryHard(b *testing.B) {
	runBenchmark(b, solver.LevelHard, solver.LevelVeryHard, 10)
}

// Initial state (bulks of 10):
// * 1	1844170500 ns/op	207972736 B/op	 1415498 allocs/op
// * Generations: 20, ~Elapsed: 184.406731ms, ~Retries: 2.100,
//   Stages: [{20 0 0} {20 0 0} {20 0 454} {20 0 0} {20 0 0} {20 0 0} {20 4 0} {16 16 0}]
// Fixes (and change to bulk of 100):
// * 1	12233012125 ns/op	1107809424 B/op	 7944966 allocs/op
// * Generations: 200, ~Elapsed: 122.328782ms, ~Retries: 1.230,
//   Stages: [{200 0 0} {200 6 0} {194 36 0} {158 158 0}]

func BenchmarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelEvil, solver.LevelEvil, 100)
}

// Initial state (bulks of 10):
// * 1	41888266875 ns/op	4805446688 B/op	32630397 allocs/op
// * Generations: 20, ~Elapsed: 4.188811889s, ~Retries: 50.000,
//   Stages: [{20 0 0} {20 0 0} {20 0 10674} {20 0 0} {20 0 0} {20 0 0} {20 6 0} {14 14 6}]
// Fixes:
// * 1	16119979791 ns/op	1435290304 B/op	10328974 allocs/op
// * Generations: 20, ~Elapsed: 1.611989s, ~Retries: 15.900,
//   Stages: [{20 0 0} {20 0 0} {20 4 0} {16 16 0}]

func BenchmarkDarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelDarkEvil, solver.LevelDarkEvil, 10)
}

func BenchmarkNightmareOrBlackHole(b *testing.B) {
	runBenchmark(b, solver.LevelNightmare, solver.LevelBlackHole, 10)
}

func runBenchmark(b *testing.B, min, max solver.Level, count int) {
	internal.Stats.Reset()
	ctx := b.Context()
	g := generators.New(&generators.Options{MinLevel: min, MaxLevel: max, Count: count})
	defer func() {
		// TODO: we can move recover to the generator itself
		msg := recover()
		if msg == nil {
			return
		}
		b.Fatalf("generator panicked with Seed: %d: %v\n%s", g.Seed(), msg, string(debug.Stack()))
	}()
	for b.Loop() {
		res := g.Generate(ctx)
		if len(res) == 0 {
			b.Fatalf("failed to generate any result")
		}
		for ri, res := range res {
			if res.Status != solver.StatusSucceeded {
				b.Fatalf("failed to generate board at result %d: %s", ri, res.Error)
			}
			if res.Steps.Level >= solver.LevelNightmare {
				b.Logf("generated[%s][%d]: %s. %s", res.Steps.Level, ri, boards.Serialize(res.Input), &res.Steps)
			}
		}
	}
	b.Log(internal.Stats.Game().String())
	// solution stats look good, no longer needed here
	// b.Log(generators.Stats.Solution().String())
}
