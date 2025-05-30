package generators_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
)

func TestGeneratorFast(t *testing.T) {
	const loop = 10
	ctx := t.Context()
	for level := solver.LevelEasy; level <= generators.FastGenerationCap; level++ {
		t.Run(level.String(), func(t *testing.T) {
			for range loop {
				g := generators.New()
				res := g.Generate(ctx, level, nil)
				assert.Assert(t, res != nil)
				assert.Equal(t, res.Status, solver.StatusSucceeded)
				assert.Check(t, res.Steps.Level >= level, "expected at least level %s, got %s", level, res.Steps.Level)
			}
		})
	}
}

// Initial state:
// * 1749	    653157 ns/op	  260671 B/op	    1535 allocs/op
// * Generations: 1749, ~Elapsed: 554.192µs, ~Retries: 1, ~Complexity: 117.87
// * Solutions: 1749, ~Elapsed: 97.204µs, ~Retries: 10.2

func BenchmarkEasy(b *testing.B) {
	runBenchmark(b, solver.LevelEasy)
}

// Initial state:
// * 519	   2266841 ns/op	  779299 B/op	    3339 allocs/op
// * Generations: 519, ~Elapsed: 2.168733ms, ~Retries: 2.45, ~Complexity: 574.43
// * Solutions: 519, ~Elapsed: 95.971µs, ~Retries: 9.6

func BenchmarkHard(b *testing.B) {
	runBenchmark(b, solver.LevelHard)
}

func BenchmarkVerHard(b *testing.B) {
	runBenchmark(b, solver.LevelVeryHard)
}

// Initial state:
// * 10	 106926608 ns/op	36740413 B/op	  149229 allocs/op
// * Generations: 10, ~Elapsed: 106.817354ms, ~Retries: 94.4, ~Complexity: 4396.20
// * Solutions: 10, ~Elapsed: 103.583µs, ~Retries: 9

func BenchmarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelEvil)
}

func BenchmarkDarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelDarkEvil)
}

func BenchmarkNightmare(b *testing.B) {
	runBenchmark(b, solver.LevelNightmare)
}

func BenchmarkBlackHole(b *testing.B) {
	runBenchmark(b, solver.LevelBlackHole)
}

func runBenchmark(b *testing.B, level solver.Level) {
	generators.Stats.Reset()
	ctx := b.Context()
	g := generators.New()
	for b.Loop() {
		res := g.Generate(ctx, level, nil)
		if res.Status != solver.StatusSucceeded {
			b.Fatalf("failed to generate board: %s", res.Error)
		}
		if res.Steps.Level >= solver.LevelNightmare {
			b.Log("generated", res.Steps.Level, ":", boards.Serialize(res.Input), &res.Steps)
		}
	}
	b.Log(generators.Stats.Game().String())
	b.Log(generators.Stats.Solution().String())
}
