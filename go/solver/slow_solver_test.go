package solver_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/solver"
)

var slowSolve = []testBoard{
	{
		// 21037 ns
		board:    "7C9A3A4J917B82B39E5A7A5C1B4A1G5D7F8C41B2B8A",
		expected: solver.StatusMoreThanOneSolution,
	},
	{
		// 70508 ns
		board:         "71A3B6B4A9A6C1A56A4A9B6A3B12E2B1D2B3A9C7A1D2B93D9A1C34A",
		expected:      solver.StatusSucceeded,
		expectedLevel: solver.LevelVeryHard,
	},
	{
		// 193476 ns
		board:         "A8C6A4A7E1A652A8D3A4A7286D816H4E6B8B1I2C4A3C",
		expected:      solver.StatusSucceeded,
		expectedLevel: solver.LevelEvil,
	},
	{
		// 12217 ns
		board:         "E87B4D398A82B691I91B49A12B7A1A2C4B73A2C2C74A3B3A18C2",
		expected:      solver.StatusSucceeded,
		expectedLevel: solver.LevelMedium,
	},
	{
		// 12322 ns
		board:    "A61F4A58D68A7A63A2C46912A8F51A1C85A6D24C16A2B8A47E76A2",
		expected: solver.StatusMoreThanOneSolution,
	},
	{
		// 24561 ns
		board:    "B32B1C1A8D27A29C65278A6C13C2958C9D23F85B25D16A9G",
		expected: solver.StatusMoreThanOneSolution,
	},
	{
		// 13682 ns
		board:    "B9C3122C38E39B4A84D692B9B8A143A1A3D6A8C12B9B6D45C9D",
		expected: solver.StatusMoreThanOneSolution,
	},
	{
		// 20395 ns
		board:    "F36A9B6C7B379C2C348A6B8947A1B22A6B9E9A4C3A4A1A79B76D2B",
		expected: solver.StatusMoreThanOneSolution,
	},
}

// start: 5876	    196460 ns/op	   42394 B/op	     318 allocs/op

func BenchmarkSlowSolve(b *testing.B) {
	benchRun(b, solver.ActionSolve, slowSolve[:])
}
