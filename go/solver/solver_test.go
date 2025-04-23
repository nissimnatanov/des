package solver_test

import (
	"encoding/json"
	"testing"

	"github.com/nissimnatanov/des/go/board"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
)

func TestSolveSanity(t *testing.T) {
	board.SetIntegrityChecks(true)

	ctx := t.Context()

	for _, sample := range sampleBoards {
		t.Run(sample.name, func(t *testing.T) {
			b, err := board.Deserialize(sample.board)
			assert.NilError(t, err)

			s := solver.New(&solver.Options{
				Action:            solver.ActionSolve,
				MaxRecursionDepth: 5,
			})

			// Solve the board
			res := s.Run(ctx, b)
			assert.NilError(t, res.Error)

			assert.Equal(t, res.Status, solver.StatusSucceeded)
			assert.Assert(t, res.Steps.Level >= solver.LevelNightmare)
			assert.Equal(t, res.Solutions.Size(), 1)
			sol := res.Solutions.At(0)
			assert.Assert(t, sol.IsValid())
			assert.Assert(t, sol.IsSolved())

			solStr := board.Serialize(sol)
			assert.Equal(t, solStr, sample.solution)

			resJSON, err := json.MarshalIndent(res, "", "  ")
			assert.NilError(t, err)
			t.Log(string(resJSON))
			t.Fail()
		})
	}
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
// max recursion = 3:           2	 617901208 ns/op	10625184 B/op	   37093 allocs/op
// max recursion = 6:           4	 308006740 ns/op	 8763860 B/op	   33233 allocs/op
// max recursion = 8:           8	 125123255 ns/op	 9652900 B/op	   38615 allocs/op
// max recursion = 8:           8	 125123255 ns/op	 9652900 B/op	   38615 allocs/op
// max recursion = 9:           10	 100261408 ns/op	11426868 B/op	   46389 allocs/op
// max recursion = 10:			12	  99804208 ns/op	13643776 B/op	   55795 allocs/op
// max recursion = 11: 			12	  99433708 ns/op	13644170 B/op	   55795 allocs/op
// max recursion = 12: 			10	 101927704 ns/op	15431537 B/op	   63405 allocs/op
// max recursion = 14:			10	 105331192 ns/op	18319923 B/op	   75798 allocs/op
// values first in base board   12	  98538781 ns/op	13643741 B/op	   55795 allocs/op
// the only choice find all instead of return on the first:
// 								13	  87288247 ns/op	13644072 B/op	   55795 allocs/op
// intro identify pairs:		30	  37168156 ns/op	 3566422 B/op	   14393 allocs/op
//
// same, change solve to prove:	19	  61008452 ns/op	 5671026 B/op	   23082 allocs/op
// trial and error indexes and board cache:
// 								19	  59378553 ns/op	  343362 B/op	   11614 allocs/op
// trial and error cache allowed+index:
// 								20	  57833294 ns/op	  530829 B/op	   17364 allocs/op
// trial and error only sort by allowed size, ignore combined value, and use slices.SortFunc:
// 								44	  27610850 ns/op	   20978 B/op	      58 allocs/op

func BenchmarkProve(b *testing.B) {
	// board.SetIntegrityChecks(true)

	ctx := b.Context()

	// Create a new board
	bd, err := board.Deserialize(sampleBoards[0].board)
	assert.NilError(b, err)

	// Create a new solver
	s := solver.New(&solver.Options{
		Action: solver.ActionProve,
	})

	for b.Loop() {
		res := s.Run(ctx, bd.Clone(board.Play))
		assert.NilError(b, res.Error)
		assert.Equal(b, res.Status, solver.StatusSucceeded)
		assert.Assert(b, res.Steps.Level >= solver.LevelNightmare)
		assert.Equal(b, res.Solutions.Size(), 1)
	}
}
