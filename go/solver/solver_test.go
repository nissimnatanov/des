package solver_test

import (
	"slices"
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/solver"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestMoreThanOneSolution(t *testing.T) {
	tb := testBoard{
		name:     "More than one solution",
		board:    "27B4A1A9B4O6A8D51A9B6D4D1B61B6D5B7A2D5C13C",
		expected: solver.StatusMoreThanOneSolution,
	}
	testSanity(t, solver.ActionProve, []testBoard{tb})
	// solver must also detect multiple solutions (it runs prove internally)
	testSanity(t, solver.ActionSolve, []testBoard{tb})
}

func TestSolveSpecificBoard(t *testing.T) {
	testSanity(t, solver.ActionSolve, benchBoards[:1])
}

func TestSolveSanity(t *testing.T) {
	testSanity(t, solver.ActionSolve, slices.Concat(benchBoards, otherBoards))
}

func TestProveSanity(t *testing.T) {
	testSanity(t, solver.ActionProve, slices.Concat(benchBoards, otherBoards))
}

func TestSolveFastSanity(t *testing.T) {
	testSanity(t, solver.ActionSolveFast, slices.Concat(benchBoards, otherBoards))
}

func TestCacheProve(t *testing.T) {
	cache := solver.NewCache()
	var firstRun bool
	var firstRunMissCount int64
	allBoards := slices.Concat(benchBoards, otherBoards)
	t.Run("all first", func(t *testing.T) {
		testSanity(t, solver.ActionProve, allBoards, cache)
		stats := cache.Stats()
		// cache hits must be 0 for the first run since recursion algo does not solve same board twice
		t.Log(stats)
		firstRun = true
		assert.Equal(t, stats.HitCount, int64(0), "expected no cache hits on proves")
		firstRunMissCount = stats.MissCount
		assert.Assert(t, firstRunMissCount > 0, "expected a lot of cache misses on first run for recursion")
	})
	t.Run("all second", func(t *testing.T) {
		if !firstRun {
			t.Skip("cannot run without the all first step")
		}
		cache.ResetStats()
		// running with cache again should not change the results
		testSanity(t, solver.ActionProve, allBoards, cache)
		stats := cache.Stats()
		t.Log(stats)
		assert.Equal(t, stats.MissCount, int64(0), "expected no cache misses on second run")
		// the cache hits must be > 0 since the cache is used, but they must be less than the misses
		// in the first run since the recursion algo stops the recursion from going deeper hence
		// saving many more recursive calls from even starting
		assert.Assert(t, stats.HitCount > 0, "expected some cache hits on second run, got %d", stats.HitCount)
		assert.Assert(t, stats.HitCount < firstRunMissCount,
			"expected less cache hits on second run than misses on first run, but got %d new hits vs %d old misses",
			stats.HitCount, firstRunMissCount)
	})
}

func TestCacheSolve(t *testing.T) {
	t.Run("al escargot", func(t *testing.T) {
		cache := solver.NewCache()
		var firstRun bool
		var firstResult *solver.Result
		var firstRunMissCount int64
		alEscargotInd := slices.IndexFunc(benchBoards, func(tb testBoard) bool {
			return tb.name == "al escargot"
		})
		assert.Assert(t, alEscargotInd >= 0, "al escargot board not found in benchBoards")
		t.Run("first", func(t *testing.T) {
			results := testSanity(t, solver.ActionSolve, benchBoards[alEscargotInd:alEscargotInd+1], cache)
			stats := cache.Stats()
			firstRun = true
			// cache hits must be 0 for the first run since recursion algo does not solve same board twice
			t.Log(stats)
			// all hits on the first run are unknown hits
			assert.Equal(t, stats.HitCount, stats.UnknownHitCount,
				"expected only unknown cache hits on first run")
			firstRunMissCount = stats.MissCount
			assert.Assert(t, firstRunMissCount > 0, "expected a lot of cache misses on first run for recursion")
			assert.Assert(t, len(results) == 1, "expected one result, got %d", len(results))
			firstResult = results[0]
		})
		t.Run("second", func(t *testing.T) {
			if !firstRun {
				t.Skip("cannot run without the all first step")
			}
			cache.ResetStats()
			// running with cache again should not change the results
			secondResults := testSanity(t, solver.ActionSolve, benchBoards[alEscargotInd:alEscargotInd+1], cache)
			stats := cache.Stats()
			t.Log(stats)
			assert.Equal(t, stats.MissCount, int64(0), "expected no cache misses on second run")
			assert.Assert(t, stats.HitCount < firstRunMissCount,
				"expected less cache hits on second run than misses on first run, but got %d new hits vs %d old misses",
				stats.HitCount, firstRunMissCount)
			if t.Failed() {
				t.Logf("first non-cache result to compare:\n %+v", firstResult)
				assert.Assert(t, len(secondResults) == 1, "expected one result, got %d", len(secondResults))
				// diff them side-by-side
				firstResString := firstResult.String()
				secondResString := secondResults[0].String()
				assert.Equal(t, firstResString, secondResString)
			}
		})
	})
	t.Run("all", func(t *testing.T) {
		cache := solver.NewCache()
		var firstRun bool
		var firstRunMissCount int64
		allBoards := slices.Concat(benchBoards, otherBoards)
		t.Run("first", func(t *testing.T) {
			testSanity(t, solver.ActionSolve, allBoards, cache)
			stats := cache.Stats()
			firstRun = true
			// cache hits must be 0 for the first run since recursion algo does not solve same board twice
			t.Log(stats)

			assert.Equal(t, stats.HitCount, stats.UnknownHitCount, "expected only unknown cache hits on first run")
			firstRunMissCount = stats.MissCount
			assert.Assert(t, firstRunMissCount > 0, "expected a lot of cache misses on first run for recursion")
		})
		t.Run("second", func(t *testing.T) {
			if !firstRun {
				t.Skip("cannot run without the all first step")
			}
			cache.ResetStats()
			// running with cache again should not change the results
			testSanity(t, solver.ActionSolve, allBoards, cache)
			stats := cache.Stats()
			t.Log(stats)
			assert.Equal(t, stats.MissCount, int64(0), "expected no cache misses on second run")
			assert.Assert(t, stats.HitCount < firstRunMissCount,
				"expected less cache hits on second run than misses on first run, but got %d new hits vs %d old misses",
				stats.HitCount, firstRunMissCount)
		})
	})
}

func testSanity(t *testing.T, action solver.Action, testBoards []testBoard, opts ...solver.Option) []*solver.Result {
	boards.SetIntegrityChecks(true)

	hasCache := false
	for _, opt := range opts {
		if _, ok := opt.(*solver.Cache); ok {
			hasCache = true
			break
		}
	}
	ctx := t.Context()
	s := solver.New()
	allResults := make([]*solver.Result, 0, len(testBoards))
	optsWithAction := slices.Concat([]solver.Option{action}, opts)
	for _, sample := range testBoards {
		t.Run(sample.name, func(t *testing.T) {
			b, err := boards.Deserialize(sample.board)
			assert.NilError(t, err)

			// Solve the board
			res := s.Run(ctx, b, optsWithAction...)
			assert.NilError(t, res.Error)
			allResults = append(allResults, res)

			expected := solver.StatusSucceeded
			if sample.expected != solver.StatusUnknown {
				expected = sample.expected
			}
			assert.Equal(t, res.Status, expected)
			if expected == solver.StatusSucceeded {
				if action == solver.ActionSolve {
					if sample.expectedLevel != solver.LevelUnknown {
						assert.Check(t, cmp.Equal(res.Level, sample.expectedLevel))
					}
					if sample.expectedComplexity > 0 {
						assert.Check(t, cmp.Equal(res.Complexity, sample.expectedComplexity))
					}
					expectedRecDepth := int8(0)
					switch {
					case res.Steps[solver.TrialAndErrorStepName][solver.StepComplexityRecursion4] > 0:
						expectedRecDepth = 4
					case res.Steps[solver.TrialAndErrorStepName][solver.StepComplexityRecursion3] > 0:
						expectedRecDepth = 3
					case res.Steps[solver.TrialAndErrorStepName][solver.StepComplexityRecursion2] > 0:
						expectedRecDepth = 2
					case res.Steps[solver.TrialAndErrorStepName][solver.StepComplexityRecursion1] > 0:
						expectedRecDepth = 1
					}
					assert.Equal(t, res.RecursionDepth, expectedRecDepth)

					if !hasCache && res.RecursionDepth != 0 {
						// solving with recursion depth set directly must not change complexity
						withMinRD := append([]solver.Option{solver.WithMinRecursionDepth(res.RecursionDepth)}, optsWithAction...)
						resWithMinRD := s.Run(ctx, b, withMinRD...)
						assert.NilError(t, resWithMinRD.Error)
						assert.Equal(t, resWithMinRD.Status, res.Status)
						assert.Check(t, cmp.Equal(resWithMinRD.Level, res.Level))
						assert.Check(t, cmp.Equal(resWithMinRD.Complexity, res.Complexity))
					}
				}

				assert.Equal(t, len(res.Solutions), 1)
				sol := res.Solutions[0]

				solStr := boards.Serialize(sol)
				assert.Check(t, cmp.Equal(solStr, sample.solution))
			}

			if sample.failToLog || t.Failed() {
				t.Log(res.String())
				t.Fail()
			}
		})
	}
	return allResults
}

// start: 			2    899252146 ns/op    840329088 B/op   3602781 allocs/op
// cloneInto
// remove remove nested ctx, time, result.options:
// use int8:
// cache nested runner:
// reset only related mask:
//
//	2	 501356520 ns/op	10721440 B/op	   39630 allocs/op
//
// max recursion = 10:
//
//	12	  99804208 ns/op	13643776 B/op	   55795 allocs/op
//
// max recursion = 14:
//
//	10	 105331192 ns/op	18319923 B/op	   75798 allocs/op
//
// values first in base board
// the only choice find all instead of return on the first:
// intro identify pairs:
//
//	30	  37168156 ns/op	 3566422 B/op	   14393 allocs/op
//
// same, change solve to prove:
//
//	19	  61008452 ns/op	 5671026 B/op	   23082 allocs/op
//
// trial and error allowed+index and board cache:
//
//	20	  57833294 ns/op	  530829 B/op	   17364 allocs/op
//
// trial and error only sort by allowed size, ignore combined value, and use slices.SortFunc:
//
//	44	  27610850 ns/op	   20978 B/op	      58 allocs/op
//
// replace value counts with free cell count only:
// allowed value cache is always valid, remove row/col/square count caches:
//
//	48	  24858786 ns/op	   19217 B/op	      58 allocs/op
//
// make Board and Solution structs instead of interfaces for performance (allows inlining):
//
//	70	  16535683 ns/op	   22024 B/op	      64 allocs/op
//
// identify all related pairs (not just first one):
// identify triplets:
//
//	72	  15989164 ns/op	   22040 B/op	      64 allocs/op
//
// single in sequence, recursion fixes, disable identify pairs and triplets in Prove:
//
//	208	   5643597 ns/op	   37840 B/op	      99 allocs/op
//
// with hint01 and bitset improvements:
//
//	222	   5334709 ns/op	   37840 B/op	      99 allocs/op
//
// continue recursion on disallowed values only:
// optimizations, bug fixes
// minor improvements in allowed:
//
//	231	   5075807 ns/op	   37984 B/op	     100 allocs/op
//
// reintroduce row/col/square value caches:
//
//	236	   5013428 ns/op	   40480 B/op	     100 allocs/op
//
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
// Improve the only choice in sequence and trial-and-error, remove object cache:
// - first only		649	   1689756 ns/op	  293009 B/op	    1109 allocs/op
// - all			416	   2803281 ns/op	  440462 B/op	    1683 allocs/op
// Improve the constraint algorithms:
// - first only		727	   1595624 ns/op	  293023 B/op	    1109 allocs/op
// - all			444	   2685913 ns/op	  440461 B/op	    1683 allocs/op
// Perf improvements:
// - first only		704	   1645071 ns/op	  423112 B/op	    2180 allocs/op
// - all			422	   2777256 ns/op	  634270 B/op	    3270 allocs/op

func BenchmarkProveFirstOnly(b *testing.B) {
	benchRun(b, solver.ActionProve, benchBoards[:1])
}

func BenchmarkProveAll(b *testing.B) {
	benchRun(b, solver.ActionProve, benchBoards)
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
// Always Prove on Solve to avoid deep recursion on boards with many solutions:
// - first only		80	  14702170 ns/op	   25258 B/op	      74 allocs/op
// - all			19	  58728298 ns/op	   76536 B/op	     250 allocs/op
// Improve the only choice in sequence and trial-and-error, remove object cache:
// - first only		103	  11445505 ns/op	  542710 B/op	    4076 allocs/op
// - all			26	  45047901 ns/op	 1389972 B/op	   12756 allocs/op
// - fast first		1384	861829 ns/op	  146616 B/op	     560 allocs/op
// - fast all		818	   1429418 ns/op	  219543 B/op	     854 allocs/op
// Improve the constraint algorithms:
// - first only		126	   9468151 ns/op	  542741 B/op	    4076 allocs/op
// - all			30	  38299142 ns/op	 1390153 B/op	   12756 allocs/op
// - fast first		1311    794930 ns/op	  146608 B/op	     560 allocs/op
// - fast all		825	   1336013 ns/op	  219549 B/op	     854 allocs/op
// Cache and bug fix in layered recursion calculations:
// - first only		123	   9647324 ns/op	  743451 B/op	    9532 allocs/op
// - all			28	  38178159 ns/op	 2055890 B/op	   30729 allocs/op
// - fast first		1441    829220 ns/op	  164023 B/op	    1092 allocs/op
// - fast all		848	   1417896 ns/op	  246973 B/op	    1659 allocs/op
// Perf improvements:
// - first only		128	   9225894 ns/op	  841715 B/op	    9532 allocs/op
// - all			32	  36221208 ns/op	 2224411 B/op	   30729 allocs/op
// - fast first		1281    821332 ns/op	  209845 B/op	    1092 allocs/op
// - fast all		860	   1386738 ns/op	  315825 B/op	    1659 allocs/op

func BenchmarkSolveFirstOnly(b *testing.B) {
	benchRun(b, solver.ActionSolve, benchBoards[:1])
}

func BenchmarkSolveAll(b *testing.B) {
	benchRun(b, solver.ActionSolve, benchBoards)
}

func BenchmarkSolveFastFirstOnly(b *testing.B) {
	benchRun(b, solver.ActionSolveFast, benchBoards[:1])
}

func BenchmarkSolveFastAll(b *testing.B) {
	benchRun(b, solver.ActionSolveFast, benchBoards)
}

func benchRun(b *testing.B, action solver.Action, testBoards []testBoard) {
	prev := solver.DisableNLog(true) // disable logging in benchmarks
	defer solver.DisableNLog(prev)

	ctx := b.Context()

	var parsed []*boards.Game
	for _, sample := range testBoards {
		bd, err := boards.Deserialize(sample.board)
		assert.NilError(b, err)
		parsed = append(parsed, bd)
	}

	// Create a new solver
	s := solver.New()

	for b.Loop() {
		for i, bd := range parsed {
			res := s.Run(ctx, bd, action)
			expected := testBoards[i].expected
			if expected == solver.StatusUnknown {
				expected = solver.StatusSucceeded
			}
			assert.NilError(b, res.Error)
			assert.Equal(b, res.Status, expected)
			if expected == solver.StatusSucceeded {
				if action.LevelRequested() {
					if testBoards[i].expectedLevel != solver.LevelUnknown {
						assert.Check(b, cmp.Equal(res.Level, testBoards[i].expectedLevel))
					} else {
						// if no expected level is set, we expect at least LevelDarkEvil
						// since the boards are not trivial
						assert.Assert(b, res.Level >= solver.LevelDarkEvil)
					}
				}
				assert.Equal(b, len(res.Solutions), 1)
			}
		}
	}
}
