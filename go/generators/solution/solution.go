package solution

import (
	"slices"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
	"github.com/nissimnatanov/des/go/internal/random"
	"github.com/nissimnatanov/des/go/internal/stats"
)

// Generate a new solution for the Sudoku board.
func Generate(r *random.Random) *boards.Solution {
	if r == nil {
		r = random.New()
	}
	g := solutionGenerator{rand: r}
	return g.generate(solutionSquareOrder)
}

type solutionGenerator struct {
	rand *random.Random
}

func (g *solutionGenerator) setSquareValues(board *boards.Game, square int, values values.Values) {
	for si, v := range values {
		board.SetReadOnly(indexes.IndexFromSquare(square, si), v)
	}
}

const maxOptimizationDepth = 4
const maxTriesForRandomSquares = 1

func (g *solutionGenerator) tryFillSeq(board *boards.Game, seq indexes.Sequence, nextSquares []int) bool {
	// pre-sort the sequence by the number of allowed values in each cell, it speeds up
	// generation time by ~30%
	slices.SortFunc(seq, func(i, j int) int {
		// sort by the number of allowed values in the square, so we fill the most constrained squares first
		return board.AllowedValues(i).Size() - board.AllowedValues(j).Size()
	})
	// since the sequence is sorted by allowed size, it is enough to check the first cell only
	if board.AllowedValues(seq[0]).Size() == 0 {
		return false
	}

	for range maxTriesForRandomSquares {
		valid := g.tryFillSequenceRec(board, seq, 0, nextSquares)
		// reassure that board is still valid (should be if we only used allowed values)
		if valid && board.IsValid() {
			return true
		}
	}

	return false
}

func (g *solutionGenerator) tryFillSequenceRec(board *boards.Game, seq indexes.Sequence, depth int, nextSquares []int) bool {
	if depth >= maxOptimizationDepth {
		valid := g.tryFillSequenceRandom(board, seq[depth:])
		if valid {
			// move on to the next square to fill, and if that fails, continue with the
			// next available values in this square
			valid = g.tryFillRemainedSquares(board, nextSquares)
		}
		if valid {
			return true
		}
		board.Reset(seq[depth:]...)
		return false
	}
	av := board.AllowedValues(seq[depth])
	var values [9]values.Value
	nv := copy(values[:], av.Values())
	random.Shuffle(g.rand, values[:nv])
	for _, v := range values[:nv] {
		board.SetReadOnly(seq[depth], v)
		if g.tryFillSequenceRec(board, seq, depth+1, nextSquares) {
			return true
		}
	}
	board.Set(seq[depth], 0)
	return false
}

func (g *solutionGenerator) tryFillSequenceRandom(board *boards.Game, seq indexes.Sequence) bool {
	for index, allowed := range board.AllowedValuesIn(seq) {
		allowedValues := allowed.Values()
		v, valid := random.Pick(g.rand, allowedValues)
		if !valid {
			// allowed values are empty, we cannot fill the reminder of this square
			board.Reset(seq...)
			return false
		}
		board.SetReadOnly(index, v)
	}
	return true
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

func (g *solutionGenerator) generate(nextSquares []int) *boards.Solution {
	start := time.Now()
	board := boards.New()
	allValues := slices.Clone(values.FullSet.Values())
	// Populate the first three squares 0, 4 (middle) and 8 (last). Since there is no intersection
	// between these squares, we can fill them with any permutation of the values.
	for soi := range 3 {
		random.Shuffle(g.rand, allValues)
		sq := nextSquares[soi]
		g.setSquareValues(board, sq, allValues)
	}
	// remove the first 3 squares from the order so that we won't use them again
	nextSquares = nextSquares[3:]

	// Fill in the values by squares, it seems to run fast enough for now. We can also try row/col or
	// other sequence variants later.
	//
	// It always succeeds after, see BenchmarkGenerateSolution for stats.
	var retries int64
retryLoop:
	for {
		if !board.IsValid() {
			panic("generated board is not valid after setting squares 0, 4, and 8")
		}

		// fill the reminder of the squares in the order defined by sqOrder
		valid := g.tryFillRemainedSquares(board, nextSquares)
		switch {
		case valid && board.IsSolved():
			// done
			break retryLoop
		case valid:
			panic("failed to fill the board, but no failed square found")
		case board.IsSolved():
			panic("board is solved, but not valid after filling squares")
		case board.FreeCellCount() != boards.Size-3*boards.SequenceSize:
			panic("some values not cleaned up after failure")
		}

		// the board is back to squares 0, 4, and 8 filled and rest of the board as empty
		retries++
	}

	sol := boards.NewSolution(board)
	// capture the stats
	elapsed := time.Since(start)
	stats.Stats.ReportOneSolution(elapsed, retries)
	return sol
}

func (g *solutionGenerator) tryFillRemainedSquares(board *boards.Game, nextSquares []int) bool {
	if len(nextSquares) == 0 {
		// for recursive calls - we are done
		return true
	}
	var squareSeq [9]int
	copy(squareSeq[:], indexes.SquareSequence(nextSquares[0]))
	return g.tryFillSeq(board, squareSeq[:], nextSquares[1:])
}
