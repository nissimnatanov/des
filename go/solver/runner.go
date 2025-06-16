package solver

import (
	"context"
	"fmt"

	"github.com/nissimnatanov/des/go/boards"
)

type runner struct {
	action Action
	input  *boards.Game

	play                  *boards.Game // the board to be solved, in Play mode
	parent                *Solver
	currentRecursionDepth int8
	maxRecursionDepth     int8
	result                *runResult
	algorithms            []Algorithm
}

func (r *runner) start(ctx context.Context) {
	if !r.Action().LevelRequested() || r.maxRecursionDepth < 2 {
		r.runWithCache(ctx)
		return
	}

	// For accurate leveling, use layer recursion: start with the recursion depth of 1 and slowly
	// increase it until we reach the max recursion or the board is solved.
	maxRecursionOriginal := r.maxRecursionDepth
	for maxRecursionCurrent := int8(1); maxRecursionCurrent <= maxRecursionOriginal; maxRecursionCurrent++ {
		r.maxRecursionDepth = maxRecursionCurrent
		r.runWithCache(ctx)
		if r.result.Status != StatusUnknown {
			return
		}
		// since we are about to increase the recursion depth, we will re-run same steps again
		if maxRecursionCurrent < maxRecursionOriginal {
			r.result.reset()
		}
	}
}

func (r *runner) runWithCache(ctx context.Context) {
	if ctx.Err() != nil {
		r.result.completeErr(ctx.Err())
		return
	}
	b := r.Board()
	if !b.IsValid() {
		r.result.complete(StatusError)
		return
	}

	// cache removed temporarily - need to check if it is beneficial
	r.runAlgos(ctx)
	return
}

func (r *runner) runAlgos(ctx context.Context) {
	b := r.Board()
	for b.FreeCellCount() > 0 {
		if ctx.Err() != nil {
			r.result.completeErr(ctx.Err())
			return
		}

		status := r.tryAlgorithms(ctx)
		if status != StatusSucceeded {
			r.result.complete(status)
			return
		}
	}

	// if we are here, we have no more free cells and the board is solved
	if !b.IsSolved() {
		// algo sets illegal value
		panic("board has no free cells but is not solved")
	}

	// algos do not report solutions
	sol := boards.NewSolution(b)
	r.result.Solutions = r.result.Solutions.Append(sol)
	if boards.GetIntegrityChecks() {
		if !boards.ContainsAll(sol, r.input) {
			panic("solution does not match the board to be solved")
		}
	}
	if len(r.result.Solutions) > 1 {
		r.result.complete(StatusMoreThanOneSolution)
	} else {
		r.result.complete(StatusSucceeded)
	}
}

func (r *runner) tryAlgorithms(ctx context.Context) Status {
	eliminationOnly := false
	b := r.Board()
	for _, algo := range r.algorithms {
		if ctx.Err() != nil {
			r.result.completeErr(ctx.Err())
			return StatusError
		}
		var startBoard *boards.Game
		if boards.GetIntegrityChecks() {
			startBoard = b.Clone(boards.Immutable)
		}

		freeBefore := b.FreeCellCount()
		status := algo.Run(ctx, r)

		if boards.GetIntegrityChecks() {
			if !boards.ContainsAll(b, startBoard) {
				panic(fmt.Errorf(
					"algo %s removed values from the board: before %v, after %v",
					algo, startBoard, b))
			}
			if !b.IsValid() {
				panic(fmt.Errorf(
					"algo %s generated failed board:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, b))
			}
		}
		if status != StatusUnknown {
			if status == StatusError && r.result.Error == nil {
				// if the algo reports an error, we want to return it
				if ctx.Err() != nil {
					r.result.completeErr(ctx.Err())
				} else {
					r.result.completeErr(fmt.Errorf("algo %s reported an error", algo))
				}
			}
			if status == StatusSucceeded &&
				b.FreeCellCount() == freeBefore &&
				!r.Action().LevelRequested() &&
				b.Hint01() < 0 {
				// If we do not need an accurate level, it is proven to be faster if we
				// try harder algorithms if the current one was only able to eliminate some
				// choices without finding a new value. The algos that eliminate only need to
				// check for the zero-or-one allowed left in the cell post elimination.
				eliminationOnly = true
				continue
			}
			return status
		}
		// with unknown status all algos must retain the board as is
		if boards.GetIntegrityChecks() {
			if !boards.Equivalent(b, startBoard) {
				panic(fmt.Errorf(
					"algo %s changed the board with unknown status:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, b))
			}
		}
	}
	if eliminationOnly {
		return StatusSucceeded
	}
	return StatusUnknown
}

func (r *runner) Action() Action {
	return r.action
}

func (r *runner) Board() *boards.Game {
	return r.play
}
func (r *runner) CurrentRecursionDepth() int {
	return int(r.currentRecursionDepth)
}
func (r *runner) MaxRecursionDepth() int {
	return int(r.maxRecursionDepth)
}
func (r *runner) AddStep(step Step, complexity StepComplexity, count int) {
	r.result.addStep(step, complexity, count)
}

func (r *runner) recursiveRun(ctx context.Context, b *boards.Game) Status {
	result := r.recursiveRunNested(ctx, b)
	// merge the sub-steps into the parent result

	// before we return, let's check the mode
	if !r.Action().LevelRequested() {
		// fast solvers do nto care about levels and can use very deep recursion
		// for efficiency, hence it does not matter how deep we are
		r.result.mergeStatsOnly(result)
		r.result.addStep(trialAndErrorStepName, StepComplexityRecursion1, 1)
		return result.Status
	}

	// Solve uses layered recursion to calculate the level, starting from maxRecursionDepth of 1
	// and increasing it until max allowed value by the options. What matters is how deep we had
	// to go to solve the board, with higher score going for the root nodes that trigger the
	// recursion (rather than the leaves).
	usedDepth := r.maxRecursionDepth - r.currentRecursionDepth

	var complexity StepComplexity
	switch usedDepth {
	case 0:
		// trialAndError only triggers nested run if current depth < max (hence delta is always at least 1)
		panic("usedDepth must be > 0, got 0")
	case 1:
		// guessed value or eliminated it with only one recursion allowed
		complexity = StepComplexityRecursion1
	case 2:
		// in order to guess this value (or eliminate it) we had to go two levels deep
		complexity = StepComplexityRecursion2
	case 3:
		// this depth is not yet reached, but hope one day we will
		complexity = StepComplexityRecursion3
		fmt.Printf("Warning: Recursion3!\n%s.\n", r.inputBoardAsString())
	case 4:
		complexity = StepComplexityRecursion4
		fmt.Printf("Warning: Recursion4!\n%s.\n", r.inputBoardAsString())
	default:
		// We should not call Solve for unproven boards, it is not effective. For proven boards,
		// we have never reached StepComplexityRecursion3, and even if we reach it one day,
		// the next one (StepComplexityRecursion4) will be even harder and likely never reached.
		// If we go beyond that, there is a terrible bug in the code somewhere.
		panic(fmt.Sprintf(
			"Unexpected usedDepth %d, something is likely wrong with the Solve algorithm.\n%s.\n",
			usedDepth, r.inputBoardAsString()))
	}
	r.result.mergeStatsOnly(result)
	r.result.addStep(trialAndErrorStepName, complexity, 1)
	return result.Status
}

func (r *runner) recursiveRunNested(ctx context.Context, b *boards.Game) *runResult {
	nested := &runner{
		action:                r.action,
		input:                 r.input,
		play:                  b,
		parent:                r.parent,
		algorithms:            r.algorithms,
		currentRecursionDepth: r.currentRecursionDepth + 1,
		maxRecursionDepth:     r.maxRecursionDepth,
		result: &runResult{
			// solutions are shared, do to clone them
			Solutions: r.result.Solutions,
			// input board is not needed, it is set below only if integrity checks are enabled
		},
	}

	if ctx.Err() != nil {
		return nested.result.completeErr(ctx.Err())
	}

	if r.currentRecursionDepth >= r.maxRecursionDepth {
		// shouldn't happen, but just in case
		return nested.result.complete(StatusUnknown)
	}

	if boards.GetIntegrityChecks() {
		// we only need input board for recursive runs if the integrity checks are enabled
		nested.input = b.Clone(boards.Immutable)
	}

	nested.runWithCache(ctx)
	return nested.result
}

// inputOrActiveBoard can be used for troubleshooting or logging abnormal cases
func (r *runner) inputBoardAsString() string {
	return "Input board: " + boards.Serialize(r.input)
}
