package solver

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/values"
)

// MaxFreeCellsForValidBoard is checked to ensure the board has enough values to start with,
// boards with less than 17 values set are mathematically proven to be illegal sudoku boards.
const MaxFreeCellsForValidBoard = boards.Size - boards.MinValidBoardSize

// maxRecursionDepthLimit is the maximum recursion depth that the solver
// can ever reach even if all other algorithms are not included and only
// the recursion is used to solve the board. For safety, set it to
// MaxFreeCellsForValidBoard, but technically speaking it would never come
// close to that number since the recursion algorithm uses 'the only choice'
// algorithm first and that would detect values.
const maxRecursionDepthLimit = 127

// Each Solver can be used on a single thread only.
type Solver struct {
	opts       Options
	algorithms []Algorithm
}

// TODO: itemize the options
func New(opts *Options) *Solver {
	if opts == nil {
		opts = &Options{}
	}
	algorithms := GetAlgorithms(opts.Action)

	return &Solver{
		opts:       *opts,
		algorithms: algorithms,
	}
}

func (s *Solver) Run(ctx context.Context, b *boards.Game) *Result {
	if b == nil {
		panic("solver.Run called with nil board")
	}

	r := &runner{
		action:                s.opts.Action,
		currentRecursionDepth: 0,
		result: Result{
			Status: StatusUnknown,
			// capture the original board AS IS, solver does not modify the Input
			Input: b,
			// before we modify the input board, we must clone it into Play mode
			Play: b.Clone(boards.Play),

			// solutions are shared so that we can deduplicate them and
			// stop once two solutions are found
			Solutions: &Solutions{},
		},
		algorithms: s.algorithms,
	}

	start := time.Now()
	defer func() {
		r.result.Elapsed = time.Since(start)
	}()

	if !b.IsValid() {
		return r.result.complete(StatusNoSolution)
	}
	if b.IsSolved() {
		// fast track the result generation for board that is already solved (e.g. Solution)
		r.result.Solutions.Add(boards.NewSolution(b))
		return r.result.complete(StatusSucceeded)
	}

	// basic checks first, do them once, there is no point in repeating them each recursion
	if b.FreeCellCount() > MaxFreeCellsForValidBoard {
		// Boards with less than 17 values are mathematically proven to be wrong.
		return r.result.complete(StatusLessThan17)
	}

	var valueCounts [boards.SequenceSize]int
	for _, v := range b.AllValues {
		if v != 0 {
			valueCounts[v-1]++
		}
	}
	missingValues := 0
	for v := values.Value(1); v <= 9; v++ {
		if valueCounts[v-1] == 0 {
			missingValues++
		}
	}
	if missingValues >= 2 {
		// There is no point to try solving boards with two or more values missing.
		return r.result.complete(StatusTwoOrMoreValuesMissing)
	}

	if s.opts.MaxRecursionDepth > maxRecursionDepthLimit ||
		s.opts.MaxRecursionDepth <= 0 {
		// without recursion it is virtually impossible to solve many boards,
		// zero is not a valid value
		// note: recursion with this package is almost 'allocation-free', and it is fast
		r.maxRecursionDepth = maxRecursionDepthLimit
	} else {
		r.maxRecursionDepth = int8(s.opts.MaxRecursionDepth)
	}

	// must run inside nested func to catch panic from run only
	// and guarantee result ref is returned
	func() {
		defer func() {
			msg := recover()
			if msg == nil {
				return
			}
			stack := string(debug.Stack())
			err, ok := msg.(error)
			if !ok {
				err = fmt.Errorf("panic: %v\n%s\n", msg, stack)
			} else {
				err = fmt.Errorf("panic: %w\n%s\n", err, stack)
			}
			r.result.completeErr(err)
		}()
		r.start(ctx)
	}()

	if r.result.Solutions.Size() > 1 && r.action.ProofRequested() && r.result.Status == StatusSucceeded {
		// if we got two solutions yet proof was requested, we must guarantee a different status
		r.result.complete(StatusMoreThanOneSolution)
	}

	// if run panics, we want to return both the error and the partial result
	return &r.result
}

type runner struct {
	action                Action
	currentRecursionDepth int8
	maxRecursionDepth     int8
	result                Result

	// to reduce calls into alloc, cache nested runner to speed up its access
	nestedCache *runner
	algorithms  []Algorithm
}

func (r *runner) start(ctx context.Context) {
	if !r.action.LevelRequested() || r.maxRecursionDepth < 2 {
		r.run(ctx)
		return
	}

	// For accurate leveling, use layer recursion: start with the recursion depth of 1 and slowly
	// increase it until we reach the max recursion or the board is solved.
	maxRecursionOriginal := r.maxRecursionDepth
	for maxRecursionCurrent := int8(1); maxRecursionCurrent <= maxRecursionOriginal; maxRecursionCurrent++ {
		r.maxRecursionDepth = maxRecursionCurrent
		r.run(ctx)
		if r.result.Status != StatusUnknown {
			return
		}
		// since we are about to increase the recursion depth, we will re-run same steps again
		if maxRecursionCurrent < maxRecursionOriginal {
			r.result.Error = nil
			r.result.Steps.reset()
		}
	}
}

func (r *runner) run(ctx context.Context) {
	if ctx.Err() != nil {
		r.result.completeErr(ctx.Err())
		return
	}
	if !r.Board().IsValid() {
		r.result.complete(StatusError)
		return
	}

	for r.Board().FreeCellCount() > 0 {
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
	if !r.Board().IsSolved() {
		// algo sets illegal value
		panic("board has no free cells but is not solved")
	}

	// algos do not report solutions
	r.result.Solutions.Add(boards.NewSolution(r.Board()))
	// we got exactly one solution, and we are done
	if boards.GetIntegrityChecks() {
		for si := range r.result.Solutions.Size() {
			sol := r.result.Solutions.At(si)
			if !boards.ContainsAll(sol, r.result.Input) {
				panic("solution does not match the board to be solved")
			}
		}
	}
	if r.result.Solutions.Size() > 1 {
		r.result.complete(StatusMoreThanOneSolution)
	} else {
		r.result.complete(StatusSucceeded)
	}
}

func (r *runner) tryAlgorithms(ctx context.Context) Status {
	eliminationOnly := false
	for _, algo := range r.algorithms {
		if ctx.Err() != nil {
			r.result.completeErr(ctx.Err())
			return StatusError
		}
		var startBoard *boards.Game
		if boards.GetIntegrityChecks() {
			startBoard = r.Board().Clone(boards.Immutable)
		}

		freeBefore := r.Board().FreeCellCount()
		status := algo.Run(ctx, r)

		if boards.GetIntegrityChecks() {
			if !boards.ContainsAll(r.Board(), startBoard) {
				panic(fmt.Errorf(
					"algo %s removed values from the board: before %v, after %v",
					algo, startBoard, r.Board()))
			}
			if !r.Board().IsValid() {
				panic(fmt.Errorf(
					"algo %s generated failed board:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, r.Board()))
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
				r.Board().FreeCellCount() == freeBefore &&
				!r.action.LevelRequested() &&
				r.Board().Hint01() < 0 {
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
			if !boards.Equivalent(r.Board(), startBoard) {
				panic(fmt.Errorf(
					"algo %s changed the board with unknown status:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, r.Board()))
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
	return r.result.Play
}
func (r *runner) CurrentRecursionDepth() int {
	return int(r.currentRecursionDepth)
}
func (r *runner) MaxRecursionDepth() int {
	return int(r.maxRecursionDepth)
}
func (r *runner) AddStep(step Step, complexity StepComplexity, count int) {
	r.result.Steps.AddStep(step, complexity, count)
}

func (r *runner) recursiveRun(ctx context.Context, b *boards.Game) *Result {
	result := r.recursiveRunNested(ctx, b)
	// merge the sub-steps into the parent result

	// before we return, let's check the mode
	if !r.action.LevelRequested() {
		// fast solvers do nto care about levels and can use very deep recursion
		// for efficiency, hence it does not matter how deep we are
		result.Steps.AddStep(trialAndErrorStepName, StepComplexityRecursion1, 1)
		r.result.Steps.Merge(&result.Steps)
		return result
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
	result.Steps.AddStep(trialAndErrorStepName, complexity, 1)
	r.result.Steps.Merge(&result.Steps)
	return result
}

func (r *runner) recursiveRunNested(ctx context.Context, b *boards.Game) *Result {
	nested := r.nestedCache
	if nested == nil {
		nested = &runner{
			action:                r.action,
			currentRecursionDepth: r.currentRecursionDepth + 1,
			maxRecursionDepth:     r.maxRecursionDepth,
			algorithms:            r.algorithms,

			// result includes ref to the shared solutions
			result: Result{
				Play:      b,
				Status:    StatusUnknown,
				Solutions: r.result.Solutions,
			},
		}
		r.nestedCache = nested
	} else {
		// we just need to set the board and reset the result's status
		nested.maxRecursionDepth = r.maxRecursionDepth
		nested.result.Play = b
		nested.result.Status = StatusUnknown
		nested.result.Error = nil
		nested.result.Steps.reset()
	}

	if ctx.Err() != nil {
		return nested.result.completeErr(ctx.Err())
	}

	if r.currentRecursionDepth >= r.maxRecursionDepth {
		// shouldn't happen, but just in case
		return nested.result.complete(StatusUnknown)
	}

	if boards.GetIntegrityChecks() {
		// we only need input board for recursive runs if the integrity checks are enabled,
		nested.result.Input = b.Clone(boards.Immutable)
	}

	nested.run(ctx)
	return &nested.result
}

// inputOrActiveBoard can be used for troubleshooting or logging abnormal cases
func (r *runner) inputBoardAsString() string {
	return "Input board: " + boards.Serialize(r.result.Input)
}
