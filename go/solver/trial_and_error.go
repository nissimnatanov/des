package solver

import (
	"context"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
)

type indexWithAllowedSize struct {
	index       int
	allowedSize int
}

// trialAndError is a recursive algorithm, please always run it after the
// theOnlyChoice algo (never as standalone)
type trialAndError struct {
	indexesCache cache[[]indexWithAllowedSize]
	testBoard    cache[*boards.Game]
}

func newTrialAndError() *trialAndError {
	return &trialAndError{
		indexesCache: cache[[]indexWithAllowedSize]{
			factory: func() []indexWithAllowedSize {
				return make([]indexWithAllowedSize, 0, MaxFreeCellsForValidBoard)
			},
			reset: func(v []indexWithAllowedSize) []indexWithAllowedSize {
				// reset the slice to be empty but keep the capacity
				return v[:0]
			},
		},
		testBoard: cache[*boards.Game]{
			factory: func() *boards.Game {
				return boards.New()
			},
		},
	}
}

func (a trialAndError) String() string {
	return "Trial and Error"
}

func (a *trialAndError) Run(ctx context.Context, state AlgorithmState) Status {
	if state.CurrentRecursionDepth() >= state.MaxRecursionDepth() {
		return StatusUnknown
	}

	b := state.Board()
	indexes := a.indexesCache.get()
	defer a.indexesCache.put(indexes) // slice is reset on next get

	for i, allowed := range b.AllAllowedValues {
		// trial and error requires at least theOnlyChoice to run first
		// hence we can safely assume allowed has at least 2 values
		indexes = append(indexes, indexWithAllowedSize{
			index:       i,
			allowedSize: allowed.Size(),
		})
	}

	// sort for faster performance (it runs ~30% faster)
	slices.SortFunc(indexes, func(tae1, tae2 indexWithAllowedSize) int {
		a1 := tae1.allowedSize
		a2 := tae2.allowedSize
		if a1 > 3 && a2 > 3 {
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
	testBoard := a.testBoard.get()
	defer a.testBoard.put(testBoard)

	var disallowedAtLeastOne bool
	for _, tei := range indexes {
		index := tei.index
		testValues := b.AllowedValues(index)
		var foundBoards []*boards.Game
		var foundUnknown bool
		var foundDisallowed bool
		for testValue := range testValues.Values {
			if ctx.Err() != nil {
				// if the context is done, we should stop
				return StatusError
			}

			b.CloneInto(boards.Play, testBoard)
			testBoard.Set(index, testValue)
			result := state.recursiveRun(ctx, testBoard)

			if result.Status == StatusUnknown {
				// remember that we had a value with unknown result and try the next one
				foundUnknown = true
				continue
			}
			if result.Status == StatusError {
				return StatusError
			}

			a.reportStep(state)
			state.MergeSteps(&result.Steps)

			if result.Status == StatusNoSolution {
				// when settings this value, the board cannot be solved, disallow it for future use
				b.DisallowValue(index, testValue)
				// we just disallowed one value, let's finish this cell since we already here
				// and restarting the recursive loop can be a waste of cycles
				// if only one value left, next loop will just set it and try to solve again
				foundDisallowed = true
				continue
			}

			// only options available are two solutions or success
			if result.Status == StatusMoreThanOneSolution {
				return StatusMoreThanOneSolution
			}

			if result.Status != StatusSucceeded {
				// should never happen
				panic("unexpected state of TrialAndError: " + result.Status.String())
			}

			if !state.Action().ProofRequested() {
				// if we do not need to prove >=1 solution, just bail report success
				// with the test board set
				copyFromTestBoard(testBoard, b)
				return StatusSucceeded
			}

			// proof requested, we must keep going with other options on the same cell to prove
			// that the found solution is the only one on that cell
			foundBoards = append(foundBoards, testBoard.Clone(boards.Immutable))
			if len(foundBoards) > 1 {
				// if we found more than one board, we can stop
				// and report that we have two solutions
				return StatusMoreThanOneSolution
			}
		}
		if len(foundBoards) == 1 {
			if !foundUnknown {
				// If proof is requested, and for each value we either found a solution or proved
				// no solution is possible, we can use the found value
				// If, however, foundUnknown is true, we cannot use the found value since we do not
				// have a hard proof that it is the only option available on the cell.
				copyFromTestBoard(foundBoards[0], b)
				return StatusSucceeded
			}
			// keep going, we have more cells to check
		}
		// no solution found, did we disallow any values?
		if foundDisallowed {
			allowed := b.AllowedValues(index)
			switch allowed.Size() {
			case 0:
				// we disallowed all the values, bail out
				return StatusNoSolution
			case 1:
				// the unknown value that is left is the only option available
				// we can set it and return
				b.Set(index, allowed.At(0))
			}

			// We disallowed one or more values on the cell, it is a bit faster (~2%) to keep going
			// with the recursive algorithm than leaving to start over. If level is not needed,
			// just continue with the next cell until we exhaust them all or set a value.
			if state.Action().LevelRequested() {
				return StatusSucceeded
			}
			disallowedAtLeastOne = true
		}
	}
	if disallowedAtLeastOne {
		// we disallowed one or more values, let's bail out and try cheaper algorithms
		return StatusSucceeded
	}

	return StatusUnknown
}

func (a *trialAndError) reportStep(state AlgorithmState) {
	var complexity StepComplexity
	switch state.CurrentRecursionDepth() {
	case 0:
		complexity = StepComplexityRecursion1
	case 1:
		complexity = StepComplexityRecursion2
	case 2:
		complexity = StepComplexityRecursion3
	case 3:
		complexity = StepComplexityRecursion4
	default:
		complexity = StepComplexityRecursion5
	}
	state.AddStep(Step(a.String()), complexity, 1)
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
