package solver

import (
	"context"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
)

const trialAndErrorStepName = Step("Trial and Error")

type indexWithAllowedSize struct {
	index  int
	weight int
}

// trialAndError is a recursive algorithm, please always run it after the
// theOnlyChoice algo (never as standalone)
type trialAndError struct {
}

func (a trialAndError) String() string {
	return string(trialAndErrorStepName)
}

func (a trialAndError) Run(ctx context.Context, state AlgorithmState) Status {
	if state.CurrentRecursionDepth() >= state.MaxRecursionDepth() {
		return StatusUnknown
	}

	b := state.Board()
	// these slice is reset on each get
	// force stack alloc
	var indexesArr [MaxFreeCellsForValidBoard]indexWithAllowedSize
	indexes := indexesArr[:0]
	// Number of allowed values is more important than the value count but only up to
	// 3 allowed values, after which the trial-and-error becomes much less effective.
	// The allowedSizeMultiplier below was chosen by testing various values and measuring
	// the performance of both prove and solve, it turns out the numbers below are the best.
	var allowedSizeMultiplier = 10
	if state.Action().ProofRequested() {
		// this multiplier works better if we have to process all the allowed values
		// and not just the first one that solves the board, thus giving more weight
		// to the allowed size performs slightly better.
		allowedSizeMultiplier = 14
	}
	const maxAllowedSizeToSort = 3
	var maxWeighToSort = maxAllowedSizeToSort * allowedSizeMultiplier

	for i, allowed := range b.AllAllowedValues {
		allowedSize := allowed.Size()
		weight := allowedSize * allowedSizeMultiplier
		if allowedSize < maxAllowedSizeToSort {
			var allValueCount int
			for _, v := range allowed.Values() {
				// adding value counts to the weight improves the trial-and-error
				// effectiveness and the overall solve speed by > ~25% and prove by ~15%
				allValueCount += b.ValueCount(v)
			}
			weight += allValueCount
		}
		// trial and error requires at least theOnlyChoice to run first
		// hence we can safely assume allowed has at least 2 values
		indexes = append(indexes, indexWithAllowedSize{
			index:  i,
			weight: weight,
		})
	}

	// sort based on the chosen weights, the lower the better
	slices.SortFunc(indexes, func(tae1, tae2 indexWithAllowedSize) int {
		a1 := tae1.weight
		a2 := tae2.weight
		if a1 > maxWeighToSort && a2 > maxWeighToSort {
			// do not bother reordering cells with more than 3 allowed, it is highly unlikely
			// this algorithm will ever need to use them and reordering them wastes ~5% of total
			// solution time
			return 0
		}
		// Important: only use the allowed size, not the 'combined' value of it,
		// adding combined value to the picture slows down the sort X 3 because
		// it needs to unnecessarily calculate the combined value and shuffle elements
		// by it even though technically speaking it is not needed at all
		return a1 - a2
	})

	// create testBoard once and reuse it
	testBoard := boards.New()

	var disallowedAtLeastOne bool
	for _, tei := range indexes {
		if ctx.Err() != nil {
			// if the context is done, we should stop the deep recursion
			return StatusError
		}
		index := tei.index
		testValues := b.AllowedValues(index).Values()
		var foundBoard *boards.Game
		var foundUnknown bool
		var foundDisallowed int
		for tvi, testValue := range testValues {
			if tvi == len(testValues)-1 && tvi == foundDisallowed {
				if !state.Action().LevelRequested() {
					// If we are at the last value and we have eliminated all others, we can skip
					// the recursion and set the value directly, this is a bit faster and also
					// more accurate in level calculation since it is equivalent to
					// the only choice in cell.
					state.AddStep(trialAndErrorStepName, StepComplexityMedium, 1)
					b.Set(index, testValue)
				} // else let the Solver detect new value and report an accurate step and level for it
				return StatusSucceeded
			}

			b.CloneInto(boards.Play, testBoard)
			testBoard.Set(index, testValue)
			resultStatus := state.recursiveRun(ctx, testBoard)
			// recursiveRun will also report the recursion step's complexity and merge the child
			// steps appropriately

			if resultStatus == StatusUnknown {
				// remember that we had a value with unknown result and try the next one
				foundUnknown = true
				continue
			}
			if resultStatus == StatusError {
				return StatusError
			}

			if resultStatus == StatusNoSolution {
				// when settings this value, the board cannot be solved, disallow it for future use
				b.DisallowValue(index, testValue)
				// we just disallowed one value, let's finish this cell since we already here
				// and restarting the recursive loop can be a waste of cycles
				// if only one value left, next loop will just set it and try to solve again
				foundDisallowed++
				continue
			}

			// only options available are two solutions or success
			if resultStatus == StatusMoreThanOneSolution {
				return StatusMoreThanOneSolution
			}

			if resultStatus != StatusSucceeded {
				// should never happen
				panic("unexpected state of TrialAndError: " + resultStatus.String())
			}

			if !state.Action().ProofRequested() && !state.Action().LevelRequested() {
				// If we do not need to prove >=1 solution, report success with the test board set.
				// Also, in a Solve mode, let's finish this cell for a stable leveling, otherwise
				// board's complexity can dramatically change after its values are shuffled.
				copyFromTestBoard(testBoard, b)
				return StatusSucceeded
			}

			// proof or accurate level requested, we must keep going with other options on the same
			// cell to prove that the found solution is the only one on that cell or to reduce the
			// 'value order bias' from the level score
			if foundBoard != nil {
				return StatusMoreThanOneSolution
			}

			foundBoard = testBoard.Clone(boards.Immutable)
		}
		// found boards can be either empty or have one board in it, we just checked above for >1
		if foundBoard != nil {
			if !foundUnknown || !state.Action().ProofRequested() {
				// if all test values except one were disallowed, and that lucky one lead to a solution,
				// then, we can conclude that this is the only solution.
				// Also, in a Solve mode, we are not asked to prove the solution, we can return early
				// to keep leveling accurate.
				copyFromTestBoard(foundBoard, b)
				return StatusSucceeded
			}
			// If proof is requested, and we have at least one inconclusive test value, we have to
			// keep going since we do not know if the board has more than one solution or not.
		}
		// no solution found, did we disallow any values?
		if foundDisallowed > 0 {
			allowed := b.AllowedValues(index).Values()
			if len(allowed) == 0 {
				// we disallowed all the values, bail out
				return StatusNoSolution
			}

			if state.Action().LevelRequested() {
				// for accurate level calculation, let the solver try simpler algorithms
				return StatusSucceeded
			}

			if len(allowed) == 1 {
				b.Set(index, allowed[0])
				state.AddStep(trialAndErrorStepName, StepComplexityMedium, 1)
			}

			// If we disallowed one or more values on the cell, or even set a value, it is a bit faster
			// (~2%) to keep going with the recursive algorithm than leaving to start over. If level is
			// not needed, just continue with the next cell until we exhaust them all.
			disallowedAtLeastOne = true
		}
	}
	if disallowedAtLeastOne {
		// we disallowed one or more values, let's bail out and try cheaper algorithms
		return StatusSucceeded
	}

	return StatusUnknown
}

func copyFromTestBoard(testBoard, b *boards.Game) {
	for i, bv := range b.AllValues {
		tv := testBoard.Get(i)
		if bv == 0 {
			if tv == 0 {
				continue
			}
			b.Set(i, tv)
			continue
		}
		if bv != tv {
			// this should never happen
			panic(
				"test board and original board have different values:\n" +
					"board\n" + b.String() + "\ntest board\n" + testBoard.String())
		}
	}
}
