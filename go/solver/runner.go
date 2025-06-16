package solver

import (
	"context"
	"fmt"

	"github.com/nissimnatanov/des/go/boards"
)

type runner struct {
	action                Action
	input                 *boards.Game
	play                  *boards.Game
	parent                *Solver
	currentRecursionDepth int8
	maxRecursionDepth     int8
	algorithms            []Algorithm
	withSteps             bool
}

func (r *runner) newRunResult() *runResult {
	rr := &runResult{}
	if r.withSteps {
		rr.Steps = make(Steps)
	}
	return rr
}

func (r *runner) start(ctx context.Context) *runResult {
	if !r.Action().LevelRequested() || r.maxRecursionDepth < 2 {
		return r.runWithCache(ctx)
	}

	// For accurate leveling, use layer recursion: start with the recursion depth of 1 and slowly
	// increase it until we reach the max recursion or the board is solved.
	var rr *runResult
	maxRecursionOriginal := r.maxRecursionDepth
	for maxRecursionCurrent := int8(1); maxRecursionCurrent <= maxRecursionOriginal; maxRecursionCurrent++ {
		r.maxRecursionDepth = maxRecursionCurrent
		rr = r.runWithCache(ctx)
		if rr.Status != StatusUnknown {
			return rr
		}
	}
	return rr
}

func (r *runner) runWithCache(ctx context.Context) *runResult {
	if ctx.Err() != nil {
		return r.newRunResult().completeErr(ctx.Err())
	}
	b := r.Board()
	if !b.IsValid() {
		return r.newRunResult().complete(StatusError)
	}

	// cache removed temporarily - need to check if it is beneficial
	return r.runAlgos(ctx)
}

func (r *runner) runAlgos(ctx context.Context) *runResult {
	runnerState := &runnerState{
		runner:    r,
		runResult: r.newRunResult(),
	}
	b := r.Board()
	for b.FreeCellCount() > 0 {
		if ctx.Err() != nil {
			return runnerState.completeErr(ctx.Err())
		}

		status := runnerState.tryAlgorithms(ctx)
		if status != StatusSucceeded {
			return runnerState.complete(status)
		}
	}

	// if we are here, we have no more free cells and the board is solved
	if !b.IsSolved() {
		// algo sets illegal value
		panic("board has no free cells but is not solved")
	}

	// algos do not report solutions
	sol := boards.NewSolution(b)
	runnerState.Solutions = runnerState.Solutions.Append(sol)
	if boards.GetIntegrityChecks() {
		if !boards.ContainsAll(sol, r.input) {
			panic("solution does not match the board to be solved")
		}
	}
	if len(runnerState.Solutions) > 1 && runnerState.Status == StatusSucceeded {
		return runnerState.complete(StatusMoreThanOneSolution)
	}
	return runnerState.complete(StatusSucceeded)
}

type runnerState struct {
	*runner
	*runResult
}

func (r *runnerState) tryAlgorithms(ctx context.Context) Status {
	eliminationOnly := false
	b := r.Board()
	for _, algo := range r.algorithms {
		if ctx.Err() != nil {
			r.completeErr(ctx.Err())
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
			if status == StatusError && r.Error == nil {
				// if the algo reports an error, we want to return it
				if ctx.Err() != nil {
					r.completeErr(ctx.Err())
				} else {
					r.completeErr(fmt.Errorf("algo %s reported an error", algo))
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
func (r *runnerState) AddStep(step Step, complexity StepComplexity, count int) {
	r.addStep(step, complexity, count)
}

func (r *runnerState) recursiveRun(ctx context.Context, b *boards.Game) Status {
	result := r.recursiveRunNested(ctx, b)
	// merge the sub-steps into the parent result
	r.mergeStatsOnly(result)

	// before we return, let's check the mode
	if !r.Action().LevelRequested() {
		// fast solvers do nto care about levels and can use very deep recursion
		// for efficiency, hence it does not matter how deep we are
		r.addStep(trialAndErrorStepName, StepComplexityRecursion1, 1)
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

	r.addStep(trialAndErrorStepName, complexity, 1)
	return result.Status
}

func (r *runner) recursiveRunNested(ctx context.Context, b *boards.Game) *runResult {
	if ctx.Err() != nil {
		return r.newRunResult().completeErr(ctx.Err())
	}

	if r.currentRecursionDepth >= r.maxRecursionDepth {
		// shouldn't happen, but just in case
		return r.newRunResult().complete(StatusUnknown)
	}
	nested := &runner{
		action:                r.action,
		input:                 r.input,
		play:                  b,
		parent:                r.parent,
		algorithms:            r.algorithms,
		currentRecursionDepth: r.currentRecursionDepth + 1,
		maxRecursionDepth:     r.maxRecursionDepth,
		withSteps:             r.withSteps,
	}
	if boards.GetIntegrityChecks() {
		// we only need input board for recursive runs if the integrity checks are enabled
		nested.input = b.Clone(boards.Immutable)
	}

	return nested.runWithCache(ctx)
}

// inputOrActiveBoard can be used for troubleshooting or logging abnormal cases
func (r *runner) inputBoardAsString() string {
	return "Input board: " + boards.Serialize(r.input)
}
