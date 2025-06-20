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
	currentRecursionDepth int8
	maxRecursionDepth     int8
	algorithms            []Algorithm
	withSteps             bool
	cache                 *Cache
}

func (r *runner) newRunResult() *runResult {
	rr := &runResult{}
	if r.withSteps {
		rr.Steps = make(Steps)
	}
	return rr
}

func (r *runner) run(ctx context.Context) *runResult {
	if !r.Action().LevelRequested() || (r.maxRecursionDepth-r.currentRecursionDepth) < 2 {
		return r.runAlgorithms(ctx, runAll)
	}

	rr := r.runAlgorithms(ctx, runNonRecursiveOnly)
	if rr.Status != StatusUnknown {
		// got the result without recursion, return it
		return rr
	}

	// For accurate leveling, use layer recursion: start with the recursion depth of current+1 and slowly
	// increase it until we reach the max recursion or the board is solved.
	maxRecursionOriginal := r.maxRecursionDepth
	for maxRecursionCurrent := r.currentRecursionDepth + 1; maxRecursionCurrent <= maxRecursionOriginal; maxRecursionCurrent++ {
		r.maxRecursionDepth = maxRecursionCurrent
		// since we already tried and applied non-recursive algos, we can start with recursion
		teResult := r.runAlgorithms(ctx, startWithRecursion)
		if teResult.Status != StatusUnknown {
			// merge the pre-recursion stats into the recursive result
			teResult.mergeStatsOnly(rr)
			return teResult
		}
	}
	return rr
}

type runAlgorithmsMode int

const (
	runAll runAlgorithmsMode = iota
	runNonRecursiveOnly
	startWithRecursion
)

func (r *runner) runAlgorithms(ctx context.Context, mode runAlgorithmsMode) *runResult {
	if ctx.Err() != nil {
		return r.newRunResult().completeErr(ctx.Err())
	}
	b := r.Board()
	if !b.IsValid() {
		return r.newRunResult().completeErr(fmt.Errorf("input board is not valid"))
	}

	var original *boards.Game
	if boards.GetIntegrityChecks() {
		// clone the original for integrity checks
		original = b.Clone(boards.Immutable)
	}
	state := &runnerState{runner: r, runResult: r.newRunResult()}
	algoState := AlgorithmState(state)
	for b.FreeCellCount() > 0 {
		if ctx.Err() != nil {
			return state.completeErr(ctx.Err())
		}
		var status Status
		if mode != startWithRecursion {
			status = state.tryNonRecursiveAlgorithms(ctx, algoState)
		} else {
			// only skip non-recursive algos once: after the first trial-and-error, go back to all algos
			mode = runAll
		}
		if status == StatusUnknown && mode != runNonRecursiveOnly {
			teResult := r.runTrialAndErrorWithCache(ctx)
			// clone the stats into current runResult
			state.mergeStatsOnly(teResult)
			status = teResult.Status
		}

		if status != StatusSucceeded {
			return state.complete(status)
		}
	}

	// if we are here, we have no more free cells and the board is solved
	if !b.IsSolved() {
		// algo sets illegal value
		panic("board has no free cells but is not solved")
	}

	// TODO: do not report same solution twice from trialAndErr, other algos do not report solutions
	sol := boards.NewSolution(b)
	state.Solutions = state.Solutions.Append(sol)
	if boards.GetIntegrityChecks() {
		if !boards.ContainsAll(sol, original) {
			panic("solution does not match the board to be solved")
		}
	}
	if len(state.Solutions) > 1 {
		return state.complete(StatusMoreThanOneSolution)
	}
	// success case
	return state.complete(StatusSucceeded)
}

type runnerState struct {
	*runner
	*runResult
}

func (r *runner) tryNonRecursiveAlgorithms(ctx context.Context, algoState AlgorithmState) Status {
	b := r.Board()
	if ctx.Err() != nil {
		return StatusError
	}
	for _, algo := range r.algorithms {
		var startBoard *boards.Game
		if boards.GetIntegrityChecks() {
			startBoard = b.Clone(boards.Immutable)
		}

		status := algo.Run(ctx, algoState)
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

	// all traditional algos failed, return back to allow caller try trial-and-error
	return StatusUnknown
}

var _ AlgorithmState = &runnerState{}

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
	if ctx.Err() != nil {
		return StatusError
	}
	status := r.recursiveRunNested(ctx, b)

	// before we return, let's check the mode
	if !r.Action().LevelRequested() {
		// prover and fast solvers do not care about level and can use very deep recursion
		// for efficiency, hence it does not matter how deep we are
		r.addStep(trialAndErrorStepName, StepComplexityRecursion1, 1)
		return status
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
	r.AddStep(trialAndErrorStepName, complexity, 1)
	return status
}

func (r *runnerState) recursiveRunNested(ctx context.Context, b *boards.Game) Status {
	if r.currentRecursionDepth >= r.maxRecursionDepth {
		return StatusUnknown
	}
	nested := &runner{
		action:                r.action,
		input:                 r.input,
		play:                  b,
		algorithms:            r.algorithms,
		currentRecursionDepth: r.currentRecursionDepth + 1,
		maxRecursionDepth:     r.maxRecursionDepth,
		cache:                 r.cache,
		withSteps:             r.withSteps,
	}

	rr := nested.run(ctx)

	// since the nested run result can be cached, we must not modify it,
	// hence, first merge the steps then report the trialAndError status only
	r.mergeStatsOnly(rr)
	return rr.Status
}

var trialAndErrorAlgo = trialAndError{}

func (r *runner) runTrialAndErrorWithCache(ctx context.Context) *runResult {
	var maxRecursionDepthTried int
	if r.action.LevelRequested() {
		maxRecursionDepthTried = int(r.maxRecursionDepth - r.currentRecursionDepth)
	}
	needIntegrityCheck := boards.GetIntegrityChecks() &&
		r.action == ActionSolve && r.cache != nil &&
		maxRecursionDepthTried == 0
	var originalPlay *boards.Game
	if needIntegrityCheck {
		originalPlay = r.play.Clone(boards.Edit)
	}

	cv, key := r.cache.get(r.Board(), r.action, maxRecursionDepthTried)
	if cv.result != nil {
		if needIntegrityCheck {
			_ = originalPlay
			tmpCache := r.cache
			r.cache = nil
			tmpPlay := r.play
			r.play = originalPlay

			// to use cache, fork a new runnerState with dedicated runResult
			tmpState := &runnerState{runner: r, runResult: r.newRunResult()}
			tmpStatus := trialAndErrorAlgo.Run(ctx, tmpState)
			if tmpState.Status != StatusUnknown {
				panic("algos should not set the status directly")
			}
			tmpState.runResult.Status = tmpStatus

			// restore prev state before checking the rest
			r.cache = tmpCache
			r.play = tmpPlay
			r.checkCacheRes(cv.result, tmpState.runResult)
		}
		if cv.result.Status != StatusUnknown {
			// we must update existing board in-place rather than replacing it (lots of code
			// captures the board as b locally)
			cv.board.CloneInto(boards.Play, r.play)
		}
		return cv.result
	}

	teState := &runnerState{runner: r, runResult: r.newRunResult()}
	teStatus := trialAndErrorAlgo.Run(ctx, teState)
	if teState.Status != StatusUnknown {
		// recursive runs only report steps and use mergeStatsOnly to update the stats
		panic("algos should not set the result status directly")
	}
	rr := teState.runResult
	rr.Status = teStatus

	// Do not cache if:
	// - the cache is not enabled (key is empty)
	// - the result is an error status
	//
	// Note: unknown (unresolved) boards are cached since we are also caching the max recursion depth with it.
	if key != NoCache && rr.Status != StatusError {
		// even if it is Prove or SolveFast, we still cache it, allowing it to be
		// reused for algos that can rely on those results, see tryCache
		// cache the result, always pass the next board (cache may ignore it for unknown)
		r.cache.set(key, r.action, CacheValue{result: rr, board: r.play}, maxRecursionDepthTried)
	}

	return rr
}

func (r *runner) checkCacheRes(cacheRes, noCacheRes *runResult) {
	// exception: cache with non-unknown result is allowed to be reused
	// in any recursion depth
	if noCacheRes.Status == StatusUnknown && cacheRes.Status != StatusUnknown {
		// allow it
		return
	}

	if noCacheRes.Status != cacheRes.Status ||
		noCacheRes.Complexity != cacheRes.Complexity ||
		noCacheRes.Count != cacheRes.Count {
		// if the result is different between cache and non-cache, it will ruin generation quality
		diff := fmt.Sprintf(
			"cache result is different from the run result, this should not happen:\nCache: %s %v\nRun: %s %v",
			cacheRes.Status, cacheRes.Steps, noCacheRes.Status, noCacheRes.Steps,
		)
		panic(diff)
	}
}

// inputOrActiveBoard can be used for troubleshooting or logging abnormal cases
func (r *runner) inputBoardAsString() string {
	return "Input board: " + boards.Serialize(r.input)
}
