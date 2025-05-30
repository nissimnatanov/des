package generators_test

import (
	"testing"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators"
	"github.com/nissimnatanov/des/go/generators/internal"
	"gotest.tools/v3/assert"
)

func TestSolutionFindFastestOrder(t *testing.T) {
	prev := boards.SetIntegrityChecks(false)
	defer boards.SetIntegrityChecks(prev)
	// Fastest so far: [0 4 8 3 5 2 1 6 7], with time: 74.272µs, retries: 10.6
	// There are more candidates, but they are more or less the same or really close.
	var order = [9]int{0, 4, 8, 3, 5, 2, 1, 6, 7}
	var fastestTimeOrder [9]int
	var fastestTimeRetries float64
	var fastestTime time.Duration = -1

	var fastestRetriesOrder [9]int
	var fastestRetries float64 = -1
	var fastestRetriesTime time.Duration

	const inLoop = 1000
	var rSeq = internal.NewRandom()

	for range 10 {
		// Reset the stats before running the benchmark
		generators.Stats.Reset()

		for range inLoop {
			var r = internal.NewRandom()
			solution := generators.GenerateSolutionWithCustomOrder(r, order[:])
			assert.Assert(t, solution != nil, "Generated solution is nil")
		}
		stats := generators.Stats.Solution()
		if fastestTime == -1 || stats.Elapsed < fastestTime {
			fastestTime = stats.Elapsed
			fastestTimeOrder = order
			fastestTimeRetries = stats.AverageRetries()
		}
		if fastestRetries == -1 || stats.AverageRetries() < fastestRetries {
			fastestRetries = stats.AverageRetries()
			fastestRetriesOrder = order
			fastestRetriesTime = stats.Elapsed
		}
		internal.RandShuffle(rSeq, order[3:])
	}
	t.Logf("Fastest time order: %v, with time: %v, retries: %v",
		fastestTimeOrder, fastestTime/time.Duration(inLoop), fastestTimeRetries)
	t.Logf("Fastest retries order: %v, with retries: %v, time: %v",
		fastestRetriesOrder, fastestRetries, fastestRetriesTime/time.Duration(inLoop))
	// t.Fail() // uncomment to fail the test and see the output
}

// GenerateSolution benchmarks the solution generation process.
// Baseline report:
// * 10508	    114152 ns/op	     768 B/op	       3 allocs/op
// * Solutions: 10508, Average Elapsed: 113.806µs, Average Retries: 11
// Pre-sort indexes in tryFillSquare:
// * 15631	     76051 ns/op	     768 B/op	       3 allocs/op
// * Solutions: 15631, Average Elapsed: 75.718µs, Average Retries: 10.9
func BenchmarkGenerateSolution(b *testing.B) {
	prev := boards.SetIntegrityChecks(false)
	defer boards.SetIntegrityChecks(prev)
	// Reset the stats before running the benchmark
	generators.Stats.Reset()
	var r = internal.NewRandom()

	for b.Loop() {
		solution := generators.GenerateSolution(r)
		assert.Assert(b, solution != nil, "Generated solution is nil")
	}

	b.Log(generators.Stats.Solution().String())
}
