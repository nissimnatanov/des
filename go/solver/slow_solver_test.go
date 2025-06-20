package solver_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/solver"
)

var slowSolve = []testBoard{
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

// start: 16888	     69800 ns/op	   43584 B/op	     272 allocs/op

func BenchmarkSlowSolve(b *testing.B) {
	benchRun(b, solver.ActionSolve, slowSolve)
}
