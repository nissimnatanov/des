package solver_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
)

/*
>>> Game Hardest 28
╔═══════╦═══════╦═══════╗
║ 6 0.0.║ 0.0.8 ║ 9 4 0.║
║ 9 0.0.║ 0.0.6 ║ 1 0.0.║
║ 0.7 0.║ 0.4 0.║ 0.0.0.║
╠═══════╬═══════╬═══════╣
║ 2 0.0.║ 6 1 0.║ 0.0.0.║
║ 0.0.0.║ 0.0.0.║ 2 0.0.║
║ 0.8 9 ║ 0.0.2 ║ 0.0.0.║
╠═══════╬═══════╬═══════╣
║ 0.0.0.║ 0.6 0.║ 0.0.5 ║
║ 0.0.0.║ 0.0.0.║ 0.3 0.║
║ 8 0.0.║ 0.0.1 ║ 6 0.0.║
╚═══════╩═══════╩═══════╝

	Sudoku Result {
	  action: Solve
	  status: Succeeded
	  level (complexity): BlackHole (26301)
	  board:  6D894A9D61C7B4D2B61J2C89B2G6C5G3A8D16B
	  elapsed (microseconds): 581739
	  value count: 22
	  steps: {
	    Easy [Single In Square, 1] X 305: 305
	    Easy [Single In Row, 1] X 17: 17
	    Easy [Single In Column, 1] X 19: 19
	    Medium [The Only Choice, 5] X 32: 160
	    Hard [Identify Pairs, 20] X 15: 300
	    Recursion1 [Trial & Error, 100] X 5: 500
	    Recursion2 [Trial & Error, 1000] X 25: 25000
	  }
	}
*/
const hardest28 = "6D894A9D61C7B4D2B61J2C89B2G6C5G3A8D16B" // complexity: 26301
const hardest28Sol = "62_5_1_7_8943_94_8_3_2_615_7_3_71_9_45_8_6_2_25_7_619_3_8_4_4_6_3_5_8_7_29_1_1_894_3_25_7_6_7_9_2_8_63_4_1_55_1_6_2_9_4_7_38_83_4_7_5_162_9_"

func TestSolveSanity(t *testing.T) {
	// board.SetIntegrityChecks(true)

	ctx := t.Context()
	//ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()

	// Create a new board
	b, err := board.Deserialize(hardest28)
	assert.NilError(t, err)

	// Create a new solver
	s := solver.New(&solver.Options{
		Action:            solver.ActionSolve,
		MaxRecursionDepth: 5,
	})

	// Solve the board
	res := s.Run(ctx, b)
	assert.NilError(t, res.Error)

	assert.Equal(t, res.Status, solver.StatusSucceeded)
	assert.Assert(t, res.StepStats.Level >= solver.LevelNightmare)
	assert.Equal(t, res.Solutions.Size(), 1)
	sol := res.Solutions.At(0)
	assert.Assert(t, sol.IsValid())
	assert.Assert(t, sol.IsSolved())

	solStr := board.Serialize(sol)
	assert.Equal(t, solStr, hardest28Sol)
}

// start: BenchmarkSolve-10		2    899252146 ns/op    840329088 B/op   3602781 allocs/op
// cloneInto:					2	 807244646 ns/op	391684240 B/op	 2901528 allocs/op
// remove nested ctx:			2	 681971688 ns/op	249680848 B/op	 1461971 allocs/op
// remove time:              	2	 650324250 ns/op	181400664 B/op	 1461959 allocs/op
// remove result.options:		2	 603979375 ns/op	147264624 B/op	 1461957 allocs/op
// remove extra result:         2	 569021250 ns/op	101755528 B/op	  750800 allocs/op
// remove extra opts, use int8:	2	 587739542 ns/op	90374112 B/op	  750795 allocs/op
// cache nested runner:         2	 535637146 ns/op	10718752 B/op	   39626 allocs/op
// reset only related mask: 	2	 501356520 ns/op	10721440 B/op	   39630 allocs/op

func BenchmarkSolve(b *testing.B) {
	// board.SetIntegrityChecks(true)

	ctx := b.Context()

	// Create a new board
	bd, err := board.Deserialize(hardest28)
	assert.NilError(b, err)

	// Create a new solver
	s := solver.New(&solver.Options{
		Action:            solver.ActionSolve,
		MaxRecursionDepth: 5,
	})

	for b.Loop() {
		res := s.Run(ctx, bd.Clone(board.Play))
		assert.NilError(b, res.Error)
		assert.Equal(b, res.Status, solver.StatusSucceeded)
		assert.Assert(b, res.StepStats.Level >= solver.LevelNightmare)
		assert.Equal(b, res.Solutions.Size(), 1)
	}
}
