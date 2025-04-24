package solver

import (
	"context"
	"slices"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/values"
)

type trialAndErrorIndex struct {
	index   int
	allowed values.Set
}

type trialAndError struct {
	indexesCache cache[[]trialAndErrorIndex]
	testBoard    cache[*boards.Play]
}

func newTrialAndError() *trialAndError {
	return &trialAndError{
		indexesCache: cache[[]trialAndErrorIndex]{
			factory: func() []trialAndErrorIndex {
				return make([]trialAndErrorIndex, 0, boards.Size-17)
			},
			reset: func(v []trialAndErrorIndex) []trialAndErrorIndex {
				// reset the slice to be empty but keep the capacity
				return v[:0]
			},
		},
		testBoard: cache[*boards.Play]{
			factory: func() *boards.Play {
				return boards.New()
			},
		},
	}
}

func (a trialAndError) String() string {
	return "Trial and Error"
}

func (a *trialAndError) Run(ctx context.Context, state AlgorithmState) Status {
	// TODO: check if we need the same layered recursion as in the C++ code
	if state.CurrentRecursionDepth() >= state.MaxRecursionDepth() {
		// this algorithm activates recursion
		return StatusUnknown
	}

	b := state.Board()
	indexes := a.indexesCache.get()
	defer a.indexesCache.put(indexes) // slice is reset on next get

	for i, allowed := range b.AllowedSets {
		if allowed.Size() == 1 {
			// if the trial algorithm is used only by itself (without other algorithms),
			// we can skip the recursion and just set the value if this is the only option
			// if we do not do this, recursion depth won't be enough to solve the board
			b.Set(i, allowed.At(0))
			theOnlyChoiceAlgo := theOnlyChoice{}
			state.AddStep(Step(theOnlyChoiceAlgo.String()), theOnlyChoiceAlgo.Complexity(), 1)
			return StatusSucceeded
		}

		indexes = append(indexes, trialAndErrorIndex{
			index:   i,
			allowed: allowed,
		})
	}

	// sort for faster performance (it runs ~30% faster)
	slices.SortFunc(indexes, func(tae1, tae2 trialAndErrorIndex) int {
		// Important: only use the allowed size, not the 'combined' value of it,
		// adding combined value to the picture slows down the sort X 3 because
		// it needs to unnecessarily calculate the combined value and shuffle elements
		// by it even though technically speaking it is not needed at all
		return tae1.allowed.Size() - tae2.allowed.Size()
	})

	// create testBoard once and reuse it
	testBoard := a.testBoard.get()
	defer a.testBoard.put(testBoard)

	for _, tei := range indexes {
		index := tei.index
		testValues := tei.allowed
		var foundBoards []*boards.Play
		foundUnknown := false
		for testValue := range testValues.Values {
			if ctx.Err() != nil {
				// if the context is done, we should stop
				return StatusError
			}

			b.CloneInto(boards.PlayMode, testBoard)
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
				b.Disallow(index, testValue)
				return StatusSucceeded
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
			foundBoards = append(foundBoards, testBoard)
			if len(foundBoards) > 1 {
				// if we found more than one board, we can stop
				// and report that we have two solutions
				return StatusMoreThanOneSolution
			}
		}
		if !foundUnknown && len(foundBoards) == 1 {
			// If proof is requested, and for each value we either found a solution or proved
			// no solution is possible, we can use the found value
			// If, however, foundUnknown is true, we cannot use the found value since we do not
			// have a hard proof that it is the only option available on the cell.
			copyFromTestBoard(foundBoards[0], b)
			return StatusSucceeded
		}
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

func copyFromTestBoard(testBoard, b *boards.Play) {
	for i := range boards.Size {
		bv := b.Get(i)
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
