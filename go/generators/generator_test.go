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

func BenchmarkEasy(b *testing.B) {
	runBenchmark(b, solver.LevelEasy)
}

func BenchmarkHard(b *testing.B) {
	runBenchmark(b, solver.LevelHard)
}

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
}
