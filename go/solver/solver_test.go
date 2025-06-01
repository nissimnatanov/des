package solver_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestMoreThanOneSolution(t *testing.T) {
	boards.SetIntegrityChecks(true)
	prover := solver.New(&solver.Options{Action: solver.ActionProve})
	ctx := t.Context()
	board, err := boards.Deserialize("27B4A1A9B4O6A8D51A9B6D4D1B61B6D5B7A2D5C13C")
	assert.NilError(t, err)

	res := prover.Run(ctx, board)
	assert.Equal(t, res.Status, solver.StatusMoreThanOneSolution)
}

func TestSolveSanity(t *testing.T) {
	testSanity(t, solver.ActionSolve)
}

func TestProveSanity(t *testing.T) {
	testSanity(t, solver.ActionProve)
}

func TestSolveFastSanity(t *testing.T) {
	testSanity(t, solver.ActionSolveFast)
}

func testSanity(t *testing.T, action solver.Action) {
	boards.SetIntegrityChecks(true)

	ctx := t.Context()
	allBoards := slices.Concat(benchBoards, otherBoards)

	for _, sample := range allBoards {
		t.Run(sample.name, func(t *testing.T) {
			b, err := boards.Deserialize(sample.board)
			assert.NilError(t, err)

			s := solver.New(&solver.Options{
				Action: action,
			})

			// Solve the board
			res := s.Run(ctx, b)
			assert.NilError(t, res.Error)

			assert.Equal(t, res.Status, solver.StatusSucceeded)
			assert.Equal(t, res.Solutions.Size(), 1)
			sol := res.Solutions.At(0)

			solStr := boards.Serialize(sol)
			assert.Check(t, cmp.Equal(solStr, sample.solution))

			resJSON, err := json.MarshalIndent(res, "", "  ")
			assert.NilError(t, err)
			t.Log(string(resJSON))
			if sample.failToLog {
				t.Fail()
			}
		})
	}
}

// start: 			2    899252146 ns/op    840329088 B/op   3602781 allocs/op
// cloneInto
// remove remove nested ctx, time, result.options:
// use int8:
// cache nested runner:
// reset only related mask:
// 				 	2	 501356520 ns/op	10721440 B/op	   39630 allocs/op
// max recursion = 10:
// 					12	  99804208 ns/op	13643776 B/op	   55795 allocs/op
// max recursion = 14:
// 					10	 105331192 ns/op	18319923 B/op	   75798 allocs/op
// values first in base board
// the only choice find all instead of return on the first:
// intro identify pairs:
// 					30	  37168156 ns/op	 3566422 B/op	   14393 allocs/op
// same, change solve to prove:
// 					19	  61008452 ns/op	 5671026 B/op	   23082 allocs/op
// trial and error allowed+index and board cache:
// 					20	  57833294 ns/op	  530829 B/op	   17364 allocs/op
// trial and error only sort by allowed size, ignore combined value, and use slices.SortFunc:
// 					44	  27610850 ns/op	   20978 B/op	      58 allocs/op
// replace value counts with free cell count only:
// allowed value cache is always valid, remove row/col/square count caches:
// 					48	  24858786 ns/op	   19217 B/op	      58 allocs/op
// make Board and Solution structs instead of interfaces for performance (allows inlining):
// 					70	  16535683 ns/op	   22024 B/op	      64 allocs/op
// identify all related pairs (not just first one):
// identify triplets:
// 					72	  15989164 ns/op	   22040 B/op	      64 allocs/op
// single in sequence, recursion fixes, disable identify pairs and triplets in Prove:
// 					208	   5643597 ns/op	   37840 B/op	      99 allocs/op
// with hint01 and bitset improvements:
// 					222	   5334709 ns/op	   37840 B/op	      99 allocs/op
// continue recursion on disallowed values only:
// optimizations, bug fixes
// minor improvements in allowed:
// 					231	   5075807 ns/op	   37984 B/op	     100 allocs/op
// reintroduce row/col/square value caches:
// 					236	   5013428 ns/op	   40480 B/op	     100 allocs/op
// update bench to solve all sample boards, not just the first one:
// - first only		232	   5047332 ns/op	   40481 B/op	     100 allocs/op
// - all			170	   6956962 ns/op	  113776 B/op	     326 allocs/op
// AllowedValuesIn, other minor improvements:
// Remove ctx checks in basic algos:
// - first only		236	   4941663 ns/op	   40480 B/op	     100 allocs/op
// - all			172	   7026087 ns/op	  113745 B/op	     326 allocs/op
// Use value counts when calculating sort weight in trial-and-error:
// - first only		367	   3138206 ns/op	   38560 B/op	      94 allocs/op
// - all			244	   4784453 ns/op	  103441 B/op	     285 allocs/op
// Remove golang enumerators and use plain slices in values.Set.Values:
// - first only		397	   3000775 ns/op	   38573 B/op	      94 allocs/op
// - all			258	   4669915 ns/op	  103462 B/op	     285 allocs/op
// Square to Row/Col and Row/Col to Square constraints:
// - first only		486	   2402109 ns/op	   35792 B/op	      90 allocs/op
// - all			277	   4184535 ns/op	  100699 B/op	     286 allocs/op
// Do not recurse on last value if others are eliminated:
// - first only		477	   2397224 ns/op	   26072 B/op	      72 allocs/op
// - all			283	   4116613 ns/op	   77955 B/op	     237 allocs/op
// Bitset improvements:
// - first only		529	   2130160 ns/op	   26082 B/op	      72 allocs/op
// - all			322	   3593113 ns/op	   77953 B/op	     237 allocs/op

func BenchmarkProveFirstOnly(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionProve,
	}, 1)
}

func BenchmarkProveAll(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionProve,
	}, len(benchBoards))
}

// start: with hint01 and bitset improvements:
// 					246	   4723883 ns/op	   33089 B/op	      82 allocs/op
// minor optimizations:
// separate the only choice in sequence out from single in sequence:
// bug fix in the only choice in sequence:
// reintroduce row/col/square value caches:
//					259	   4546273 ns/op	   35072 B/op	      83 allocs/op
// update bench to solve all sample boards, not just the first one:
// - first only		260	   4546737 ns/op	   35092 B/op	      83 allocs/op
// - all			207	   5725250 ns/op	   82961 B/op	     238 allocs/op
// AllowedValuesIn, other minor improvements:
// Remove ctx checks in basic algos:
// Use value counts when calculating sort weight in trial-and-error:
// Remove golang enumerators and use plain slices in values.Set.Values:
// Square to Row/Col constraints:
// Row/Col to Square constraints (increases back but needed for accurate level):
// Do not recurse on last value if others are eliminated:
// Bitset improvements:
// - first only		484	   2342197 ns/op	   20960 B/op	      58 allocs/op
// - all			387	   3038964 ns/op	   57329 B/op	     185 allocs/op
// - fast first		984	   1096763 ns/op	   22832 B/op	      62 allocs/op
// - fast all		614	   1824448 ns/op	   67961 B/op	     203 allocs/op
// Layered recursion and reduce value bias in scoring:
// - first only		92	  12721382 ns/op	    6896 B/op	      30 allocs/op
// - all			21	  55324716 ns/op	   26824 B/op	     117 allocs/op
// - fast first		943	   1130329 ns/op	   22960 B/op	      61 allocs/op
// - fast all		637	   1830455 ns/op	   68280 B/op	     199 allocs/op

func BenchmarkSolveFirstOnly(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionSolve,
	}, 1)
}

func BenchmarkSolveAll(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionSolve,
	}, len(benchBoards))
}

func BenchmarkSolveFastFirstOnly(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionSolveFast,
	}, 1)
}

func BenchmarkSolveFastAll(b *testing.B) {
	benchRun(b, &solver.Options{
		Action: solver.ActionSolveFast,
	}, len(benchBoards))
}

func benchRun(b *testing.B, opts *solver.Options, numBoards int) {
	ctx := b.Context()

	var parsed []*boards.Game
	for _, sample := range benchBoards {
		bd, err := boards.Deserialize(sample.board)
		assert.NilError(b, err)
		parsed = append(parsed, bd)
		numBoards--
		if numBoards <= 0 {
			break
		}
	}

	// Create a new solver
	s := solver.New(opts)

	for b.Loop() {
		for _, bd := range parsed {
			res := s.Run(ctx, bd.Clone(boards.Play))
			assert.NilError(b, res.Error)
			assert.Equal(b, res.Status, solver.StatusSucceeded)
			if opts.Action.LevelRequested() {
				assert.Assert(b, res.Steps.Level >= solver.LevelDarkEvil)
			}
			assert.Equal(b, res.Solutions.Size(), 1)
		}
	}
}
