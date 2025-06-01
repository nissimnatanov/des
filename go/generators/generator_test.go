package generators_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestGeneratorFast(t *testing.T) {
	const loop = 10
	ctx := t.Context()
	for level := solver.LevelEasy; level <= generators.FastGenerationCap; level++ {
		t.Run(level.String(), func(t *testing.T) {
			for range loop {
				g := generators.New(&generators.Options{MinLevel: level, MaxLevel: level})
				rs := g.Generate(ctx)
				assert.Check(t, cmp.Len(rs, 1), "expected exactly one result for level %s, got %d", level, len(rs))
				for _, res := range rs {
					assert.Assert(t, res != nil)
					assert.Equal(t, res.Status, solver.StatusSucceeded)
					assert.Check(t, res.Steps.Level == level, "expected at least level %s, got %s", level, res.Steps.Level)
				}
			}
		})
	}
}

// Initial state:
// * 1749	    653157 ns/op	  260671 B/op	    1535 allocs/op
// * Generations: 1749, ~Elapsed: 554.192µs, ~Retries: 1, ~Complexity: 117.87
// * Solutions: 1749, ~Elapsed: 97.204µs, ~Retries: 10.2

func BenchmarkEasyOrMedium(b *testing.B) {
	runBenchmark(b, solver.LevelEasy, solver.LevelMedium, 100)
}

func BenchmarkHardOrVeryHard(b *testing.B) {
	runBenchmark(b, solver.LevelHard, solver.LevelVeryHard, 20)
}

// Initial state:
// * 10	 106926608 ns/op	36740413 B/op	  149229 allocs/op
// * Generations: 10, ~Elapsed: 106.817354ms, ~Retries: 94.4, ~Complexity: 4396.20
// * Solutions: 10, ~Elapsed: 103.583µs, ~Retries: 9

func BenchmarkEviOrDarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelEvil, solver.LevelDarkEvil, 10)
}

func BenchmarkNightmareOrBlackHole(b *testing.B) {
	runBenchmark(b, solver.LevelNightmare, solver.LevelBlackHole, 1)
}

func runBenchmark(b *testing.B, min, max solver.Level, count int) {
	generators.Stats.Reset()
	ctx := b.Context()
	g := generators.New(&generators.Options{MinLevel: min, MaxLevel: max, Count: count})
	defer func() {
		// TODO: we can move recover to the generator itself
		msg := recover()
		if msg == nil {
			return
		}
		b.Fatalf("generator panicked with Seed: %d: %v", g.Seed(), msg)
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
	b.Log(generators.Stats.Game().String())
	// solution stats look good, no longer needed here
	// b.Log(generators.Stats.Solution().String())
}
