package generators

import (
	"slices"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

func GenerateSolution(r *Random) *boards.Solution {
	if r == nil {
		r = NewRandom()
	}
	return solutionGenerator{rand: r}.generate(solutionSquareOrder)
}

type solutionGenerator struct {
	rand *Random
}

func (g solutionGenerator) setSquareValuesReadOnly(board *boards.Game, square int, values values.Values) {
	for si, v := range values {
		board.SetReadOnly(indexes.IndexFromSquare(square, si), v)
	}
}

func (g solutionGenerator) tryFillSeq(board *boards.Game, seq indexes.Sequence) bool {
	// pre-sort the sequence by the number of allowed values in each cell, it speeds up
	// generation time by ~30%
	slices.SortFunc(seq, func(i, j int) int {
		// sort by the number of allowed values in the square, so we fill the most constrained squares first
		return board.AllowedValues(i).Size() - board.AllowedValues(j).Size()
	})
	for _, allowed := range board.AllowedValuesIn(seq) {
		if allowed.Size() == 0 {
			return false
		}
	}

	// Most of the time the sequence is filled with fewer tries, yet to avoid infinite loops let's
	// cap it to 10 tries. Benchmark shows that numbers after 10 do not improve the success rate.
	const maxTries = 10
	for range maxTries {
		var valid bool
		var v values.Value

		for index, allowed := range board.AllowedValuesIn(seq) {
			allowedValues := allowed.Values()
			v, valid = RandPick(g.rand, allowedValues)
			if valid {
				board.Set(index, v)
			}
			if !valid {
				break
			}
		}
		// reassure that board is still valid (should be if we only used allowed values)
		valid = valid && board.IsValid()
		if valid {
			return true // successfully filled the square
		}

		// reset the square and try again if we have retries left
		for _, index := range seq {
			board.Set(index, 0)
		}
	}

	return false
}

// solutionSquareOrder defines the order of squares to fill
//
// The first 3 squares can be pre-filled with any permutation of values.
// The order for others apparently has impact too, and benchmarks show that the
// order below is the fastest so far.
//
// Note: choosing the wrong order can significantly impact the performance of the solution generation,
// increasing number of retries from ~10 to over 70. Play around with the TestSolutionFindFastestOrder
// test to reset to the best order.
var solutionSquareOrder = []int{0, 4, 8, 3, 5, 2, 1, 6, 7}

func (g solutionGenerator) generate(sqOrder []int) *boards.Solution {
	start := time.Now()
	board := boards.New()
	allValues := slices.Clone(values.FullSet.Values())
	// Populate the first three squares 0, 4 (middle) and 8 (last). Since there is no intersection
	// between these squares, we can fill them with any permutation of the values.
	for soi := range 3 {
		RandShuffle(g.rand, allValues)
		sq := sqOrder[soi]
		g.setSquareValuesReadOnly(board, sq, allValues)
	}

	// avoid slice allocations in a heavy loop
	var squareSeqCache [9]int

	// Fill in the values by squares, it seems to run fast enough for now. We can also try row/col or
	// other sequence variants later.
	//
	// It always succeeds after, see BenchmarkGenerateSolution for stats.
	var retries int64
	for {
		if !board.IsValid() {
			panic("generated board is not valid after setting squares 0, 4, and 8")
		}

		// fill the reminder of the squares in the order defined by sqOrder
		for soi := 3; soi < 9; soi++ {
			sq := sqOrder[soi]
			copy(squareSeqCache[:], indexes.SquareSequence(sq))
			if !g.tryFillSeq(board, squareSeqCache[:]) {
				break
			}
		}
		// if we filled all squares and the board remained valid, we are done
		if board.IsSolved() {
			break
		}

		// restart back to squares 0, 4, and 8 filled and rest of the board as empty
		board.Restart()
		retries++
	}
	// mark all read-only
	for i := range boards.Size {
		board.SetReadOnly(i, board.Get(i))
	}
	sol := boards.NewSolution(board)
	// capture the stats
	elapsed := time.Since(start)
	Stats.reportOneSolution(elapsed, retries)
	return sol
}
