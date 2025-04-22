package solver

import (
	"context"
	"sort"

	"github.com/nissimnatanov/des/go/board"
)

type trialAndError struct {
}

func (a trialAndError) String() string {
	return "Trial and Error"
}

func (a trialAndError) Run(ctx context.Context, state AlgorithmState) Status {
	// TODO: check if we need the same layered recursion as in the C++ code
	if state.CurrentRecursionDepth() >= state.MaxRecursionDepth() {
		// this algorithm activates recursion
		return StatusUnknown
	}

	b := state.Board()
	indexes := make([]int, 0, b.FreeCellCount())
	for i := range board.Size {
		if b.IsEmpty(i) {
			indexes = append(indexes, i)
		}
		if b.AllowedSet(i).Size() == 1 {
			// if this algorithm is used only by itself (without other algorithms),
			// we can skip the recursion and just set the value if this is the only option
			// if we do not do this, recursion depth won't be enough to solve the board
			b.Set(i, b.AllowedSet(i).At(0))
			theOnlyChoiceAlgo := theOnlyChoice{}
			state.AddStep(Step(theOnlyChoiceAlgo.String()), theOnlyChoiceAlgo.Complexity(), 1)
			return StatusSucceeded
		}
	}

	//	if !state.Action().LevelRequested() {
	// sort for faster performance (it runs ~twice faster)
	sort.Slice(indexes, func(i, j int) bool {
		allowed1 := b.AllowedSet(indexes[i])
		allowed2 := b.AllowedSet(indexes[j])
		if allowed1.Size() != allowed2.Size() {
			// the less allowed the better
			return allowed1.Size() < allowed2.Size()
		}
		return allowed1.Combined() < allowed2.Combined()
	})

	// create testBoard once and reuse it
	testBoard := board.New()

	for _, index := range indexes {
		testValues := b.AllowedSet(index)
		var foundBoards []board.Board
		foundUnknown := false
		for ti := range testValues.Size() {
			if ctx.Err() != nil {
				// if the context is done, we should stop
				return StatusError
			}
			testValue := testValues.At(ti)

			b.CloneInto(board.Play, testBoard)
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
			state.MergeSteps(&result.StepStats)

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
				copyTestBoard(testBoard, b)
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
			copyTestBoard(foundBoards[0], b)
			return StatusSucceeded
		}
	}

	return StatusUnknown
}

func (a trialAndError) reportStep(state AlgorithmState) {
	var complexity StepComplexity
	switch state.CurrentRecursionDepth() {
	case 0:
		complexity = StepComplexityRecursion1
	case 1:
		complexity = StepComplexityRecursion2
	case 2:
		complexity = StepComplexityRecursion3
	default:
		complexity = StepComplexityRecursion4
	}
	state.AddStep(Step(a.String()), complexity, 1)
}

func copyTestBoard(testBoard, b board.Board) {
	for i := range board.Size {
		if b.IsEmpty(i) {
			if testBoard.IsEmpty(i) {
				continue
			}
			b.Set(i, testBoard.Get(i))
			continue
		}
		if b.Get(i) != testBoard.Get(i) {
			// this should never happen
			panic(
				"test board and original board have different values:\n" +
					"board\n" + b.String() + "\ntest board\n" + testBoard.String())
		}
	}
}
