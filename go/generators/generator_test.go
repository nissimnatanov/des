package generators_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"testing"
	"time"

	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/generators/internal"
	"github.com/nissimnatanov/des/go/internal/stats"
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
					assert.Check(t, res.Level == level, "expected level %s, got %s", level, res.Level)
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
			assert.Check(t, res.Level >= min && res.Level <= max,
				"expected level between %s and %s, got %s", min, max, res.Level)
		}
	}
}

// Initial state (bulks of 100):
// * 18	  57867785 ns/op	 7639343 B/op	   50437 allocs/op
// * Generations: 1800, ~Elapsed: 577.917µs, ~Retries: 1.000,
//   Stages: [{1800 452 0} {1348 1048 0} {300 300 0}]
// Improve the only choice in sequence and trial-and-error:
// * 26	  44879917 ns/op	13663641 B/op	   68303 allocs/op
// * Generations: 2600, ~Elapsed: 448.108µs, ~Retries: 1.000,
//   Stages: [{2600 645 0} {1955 1575 0} {380 380 0}]
// Cache and bug fix in layered recursion calculations:
// * 20	  55638904 ns/op	19075226 B/op	  147718 allocs/op
// * Generations: 2000, ~Elapsed: 552.749µs, ~Retries: 1.000,
//   Stages: [{2000 426 0} {1574 1222 0} {352 352 0}]

func BenchmarkEasyOrMedium(b *testing.B) {
	runBenchmark(b, solver.LevelEasy, solver.LevelMedium, 100)
}

// Initial state (bulks of 10):
// * 22	  59968347 ns/op	 6069760 B/op	   41416 allocs/op
// * Generations: 220, ~Elapsed: 5.992518ms, ~Retries: 2.205,
//   Stages: [{220 0 0} {220 42 0} {178 178 6}]
// Improve the only choice in sequence and trial-and-error:
// * 27	  49116306 ns/op	13561052 B/op	   63190 allocs/op
// * Generations: 270, ~Elapsed: 4.90605ms, ~Retries: 2.367,
//   Stages: [{270 0 0} {270 52 3} {218 218 6}]
// Cache and bug fix in layered recursion calculations:
// * 20	  53990679 ns/op	17271384 B/op	  136893 allocs/op
// * Generations: 200, ~Elapsed: 5.391663ms, ~Retries: 2.305,
//   Stages: [{200 1 0} {199 32 2} {167 167 10}]

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
// Improve the only choice in sequence and trial-and-error, remove object cache:
// * 1	13000732166 ns/op	3166094672 B/op	16387331 allocs/op
// * Generations: 204, ~Elapsed: 127.456325ms, ~Retries: 1.569,
//   Stages: [{226 0 0} {226 20 0} {206 28 0} {178 178 0}]
// Improve the constraint algorithms:
// * 1	10674286042 ns/op	2678858688 B/op	13787793 allocs/op
// * Generations: 204, ~Elapsed: 104.648522ms, ~Retries: 1.314,
//   Stages: [{228 0 0} {228 14 0} {214 40 0} {174 174 0}]
// Cache and bug fix in layered recursion calculations, perf improvements:
// * 1	10149820459 ns/op	4704790160 B/op	27352226 allocs/op
// * Generations: 100, ~Elapsed: 101.498177ms, ~Retries: 1.100,
//   Stages: [{115 0 0} {115 3 0} {112 24 0} {88 88 0}]
// * Solver Cache: hits=183488 (12.45%), unknown hits=1714 (0.12%), misses=1288829, sets=1288861
// Generation layering and other improvements:
// * 1	3800123250 ns/op	1583798840 B/op	 9499364 allocs/op
// * Generations: 102, ~Elapsed: 37.256077ms, ~Retries: 0.098, Stages:
//   [{102 0 0} {102 3 0} {99 99 0}]
// * Solver Cache: hits=120559 (23.64%), unknown hits=1376 (0.27%), misses=387987,
//   sets=388001, unknown sets=4221 (1.09%)
// Tune generation parameters:
// * 1	1047977750 ns/op	416294920 B/op	 2482588 allocs/op
// * Solutions: 1, ~Elapsed: 73.792µs, ~Retries: 0.000
// * Solver Cache: hits=37911 (29.59%, unknown=0.37%), misses=89727, sets=89734 (unknown=1581)
// * Game Generations: 106, ~Elapsed: 9.500251ms, ~Retries: 0.047. Stages: [
// *   [3] Total: 114, Success: 20(17.54%), Failed: 0(0.00%)
// *   [4] Total: 94, Success: 39(41.49%), Failed: 0(0.00%)
// *   [5] Total: 55, Success: 55(100.00%), Failed: 0(0.00%)
// * ], Complexities: 166 | 149.5; 1205 | 690.0; 4334 | 1973.4; 8699 | 4382.6; 9670 | 5496.1; 9670 | 5157.5
// TopN, ProveOnly and more perf improvements:
// * 3	 414796222 ns/op	136639650 B/op	  808240 allocs/op
// * Solutions: 7, ~Elapsed: 53.696µs, ~Retries: 0.000
// * Solver Cache: hits=99550 (63.55%, unknown=0.54%), misses=56242, sets=56250 (unknown=2104)
// * Game Generations: 300, ~Elapsed: 4.147318ms, ~Retries: 0.023. Stages: [
// *   [2] Total: 382, Success: 90(23.56%), Failed: 0(0.00%), ~Candidates: 94.3
// *   [3] Total: 292, Success: 195(66.78%), Failed: 0(0.00%), ~Candidates: 50.0
// *   [4] Total: 97, Success: 80(82.47%), Failed: 0(0.00%), ~Candidates: 40.0
// *   [5] Total: 17, Success: 15(88.24%), Failed: 2(11.76%), ~Candidates: 20.2
// * ], Complexities: [0] 525/275 [1] 1015/637 [2] 8158/4297 [3] 8178/5554 [4] 7903/5705 [5] 7903/7498

func BenchmarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelEvil, solver.LevelEvil, 100)
}

// Initial state (bulks of 10):
//   - 1	41888266875 ns/op	4805446688 B/op	32630397 allocs/op
//   - Generations: 20, ~Elapsed: 4.188811889s, ~Retries: 50.000,
//     Stages: [{20 0 0} {20 0 0} {20 0 10674} {20 0 0} {20 0 0} {20 0 0} {20 6 0} {14 14 6}]
//
// Fixes:
//   - 1	16119979791 ns/op	1435290304 B/op	10328974 allocs/op
//   - Generations: 20, ~Elapsed: 1.611989s, ~Retries: 15.900,
//     Stages: [{20 0 0} {20 0 0} {20 4 0} {16 16 0}]
//
// Improve the only choice in sequence and trial-and-error:
//   - 1	14836185834 ns/op	3555023984 B/op	17091575 allocs/op
//   - Generations: 20, ~Elapsed: 1.483607195s, ~Retries: 17.800,
//     Stages: [{20 0 0} {20 2 0} {18 0 0} {18 18 0}]
//
// Cache and bug fix in layered recursion calculations:
//   - 1	28215436834 ns/op	11097599904 B/op	89392167 allocs/op
//   - Generations: 10, ~Elapsed: 2.821543495s, ~Retries: 30.200,
//     Stages: [{10 0 0} {10 1 0} {9 0 0} {9 9 0}]
//   - Solver Cache: hits=3018 (10.9%), unknown hits=100 (0.4%), misses=24487, sets=24489
//
// Generation layering and other improvements:
//   - 1	16818874958 ns/op	6950772096 B/op	41741684 allocs/op
//   - Generations: 10, ~Elapsed: 1.6818873s, ~Retries: 4.000,
//     Stages: [{10 0 0} {10 0 0} {10 10 0}]
//   - Solver Cache: hits=534845 (23.61%), unknown hits=4985 (0.22%), misses=1725744,
//     sets=1725850, unknown sets=18245
//
// Tune generation parameters:
// * 1	12090348292 ns/op	4873126480 B/op	28898735 allocs/op
// * Solutions: 3, ~Elapsed: 40.069µs, ~Retries: 0.000
// * Solver Cache: hits=458344 (29.78%, unknown=0.50%), misses=1072972, sets=1073128 (unknown=21607)
// * Game Generations: 25, ~Elapsed: 483.609035ms, ~Retries: 2.280. Stages: [
// *   [3] Total: 80, Success: 3(3.75%), Failed: 0(0.00%)
// *   [4] Total: 77, Success: 15(19.48%), Failed: 0(0.00%)
// *   [5] Total: 62, Success: 11(17.74%), Failed: 51(82.26%)
// * ], Complexities: 443 | 184.0; 900 | 619.4; 2750 | 1691.1; 29082 | 8323.8; 29087 | 9585.3; 29087 | 9101.6
// TopN, ProveOnly and more perf improvements:
// * 1	7007583583 ns/op	2734085904 B/op	15746173 allocs/op
// * Solutions: 15, ~Elapsed: 49.008µs, ~Retries: 0.000
// * Solver Cache: hits=564743 (57.52%, unknown=0.37%), misses=413576, sets=413622 (unknown=9973)
// * Game Generations: 100, ~Elapsed: 70.068484ms, ~Retries: 0.450. Stages: [
// *   [3] Total: 149, Success: 17(11.41%), Failed: 0(0.00%), ~Candidates: 50.0
// *   [4] Total: 132, Success: 79(59.85%), Failed: 0(0.00%), ~Candidates: 40.0
// *   [5] Total: 53, Success: 16(30.19%), Failed: 37(69.81%), ~Candidates: 23.4
// * ], Complexities: [0] 420/236 [1] 1290/732 [2] 4395/2128 [3] 22183/6119 [4] 22183/6581 [5] 22183/13946

func BenchmarkDarkEvil(b *testing.B) {
	runBenchmark(b, solver.LevelDarkEvil, solver.LevelDarkEvil, 100)
}

func BenchmarkNightmareOrBlackHole(b *testing.B) {
	runBenchmark(b, solver.LevelNightmare, solver.LevelBlackHole, 10)
}

func runBenchmark(b *testing.B, min, max solver.Level, count int) {
	stats.Stats.Reset()
	ctx := b.Context()
	ctx, sygCancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer sygCancel()

	// benchmark logger truncates output after 10 lines, open a tmp log file instead
	// and immediately show it in vscode
	logFile := filepath.Join(os.TempDir(),
		fmt.Sprintf("bench_%s.log", time.Now().Format("20060102_150405")))

	b.Log("Log File: ", logFile)

	r := stats.Reporter{
		SkipOnSilence: 1,
		Duration:      time.Second * 10, // report every 10 seconds
		OutputFile:    logFile,
	}

	r.Run()
	defer r.Stop()

	cmd := exec.Command("code", logFile)
	err := cmd.Run()
	assert.Check(b, err, "failed to open log file with vscode")

	// force stop to print stats after 8 minutes
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Minute*8)
	defer cancel()

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
		if len(res) == 0 && ctx.Err() == nil {
			assert.Check(b, false, "failed to generate any result")
		}
		for ri, res := range res {
			if res.Status != solver.StatusSucceeded {
				if res.Status != solver.StatusError || errors.Is(res.Error, context.Canceled) {
					assert.Check(b, false,
						"failed to generate board at result %d: %s", ri, res.Error)
				}
			}
		}
	}
	// log the slow boards
	if internal.SlowBoards.HasLog() {
		slowBoards := internal.SlowBoards.Log()
		r.LogNow(slowBoards)
	}
	b.Log("Done")
}
