package solver

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/nissimnatanov/des/go/board"
	"github.com/nissimnatanov/des/go/board/values"
)

// MaxFreeCellsForValidBoard is checked to ensure the board has enough values to start with,
// boards with less than 17 values set are mathematically proven to be illegal sudoku boards.
const MaxFreeCellsForValidBoard = 64

type Solver struct {
	opts Options
}

// TODO: itemize the options
func New(opts *Options) *Solver {
	if opts == nil {
		opts = &Options{}
	}

	return &Solver{
		opts: *opts,
	}
}

func (s *Solver) Run(ctx context.Context, b board.Board) *Result {
	result := &Result{
		Status:    StatusUnknown,
		Solutions: &Solutions{},
	}

	if b == nil {
		return result.completeErr(fmt.Errorf("board is nil"))
	}

	// basic checks first, do them once, there is no point in repeating them each recursion
	if b.FreeCellCount() > MaxFreeCellsForValidBoard {
		// Boards with less than 17 values are mathematically proven to be wrong.
		return result.complete(StatusLessThan17)
	}

	missingValues := 0
	for v := values.Value(1); v <= 9; v++ {
		if b.Count(v) == 0 {
			missingValues++
		}
	}
	if missingValues >= 2 {
		// There is no point to try solving boards with two or more values missing.
		return result.complete(StatusTwoOrMoreValuesMissing)
	}

	if !b.IsValid() {
		return result.complete(StatusNoSolution)
	}

	r := &runner{
		board:                 b,
		action:                s.opts.Action,
		currentRecursionDepth: 0,
		result: Result{
			Status: StatusUnknown,

			// solutions are shared so that we can deduplicate them and
			// stop once two solutions are found
			Solutions: result.Solutions,
		},
		algorithms: GetAlgorithms(s.opts.Action),
	}
	if board.GetIntegrityChecks() {
		// capture the original board for integrity checks to make sure algos do not corrupt
		// the board and their solutions solve the input
		r.inputBoard = b.Clone(board.Immutable)
	}

	switch {
	case s.opts.MaxRecursionDepth > 32:
		// 32 is way too deep, best perf achieved around 10-15 and from there usually
		// the-only-choice algo completes the board, anything above that is useless
		r.maxRecursionDepth = 32
	case s.opts.MaxRecursionDepth <= 0:
		// without recursion it is virtually impossible to solve many boards,
		// zero is not a valid value
		r.maxRecursionDepth = 10
	default:
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
		r.run(ctx)
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
	board                 board.Board
	currentRecursionDepth int8
	maxRecursionDepth     int8
	result                Result

	// inputBoard for integrity checks of the solutions
	inputBoard board.Board

	// to reduce calls into alloc, cache nested runner to speed up its access
	nestedCache *runner
	algorithms  []Algorithm
}

func (r *runner) run(ctx context.Context) *Result {
	if ctx.Err() != nil {
		return r.result.completeErr(ctx.Err())
	}
	if !r.board.IsValid() {
		return r.result.complete(StatusError)
	}

	for r.board.FreeCellCount() > 0 {
		if ctx.Err() != nil {
			return r.result.completeErr(ctx.Err())
		}

		status := r.tryAlgorithms(ctx)
		if status != StatusSucceeded {
			return r.result.complete(status)
		}
	}

	// if we are here, we have no more free cells and the board is solved
	if !r.board.IsSolved() {
		// algo sets illegal value
		panic("board has no free cells but is not solved")
	}

	// algos do not report solutions
	r.result.Solutions.Add(board.NewSolution(r.board))
	// we got exactly one solution, and we are done
	if board.GetIntegrityChecks() {
		for si := range r.result.Solutions.Size() {
			sol := r.result.Solutions.At(si)
			if !board.ContainsAll(sol, r.inputBoard) {
				panic("solution does not match the board to be solved")
			}
			if !sol.IsSolved() {
				panic("solution is not solved")
			}
		}
	}
	if r.result.Solutions.Size() > 1 {
		return r.result.complete(StatusMoreThanOneSolution)
	}
	return r.result.complete(StatusSucceeded)
}

func (r *runner) tryAlgorithms(ctx context.Context) Status {
	for _, algo := range r.algorithms {
		if ctx.Err() != nil {
			r.result.completeErr(ctx.Err())
			return StatusError
		}
		var startBoard board.Board
		if board.GetIntegrityChecks() {
			startBoard = r.board.Clone(board.Immutable)
		}

		status := algo.Run(ctx, r)

		if board.GetIntegrityChecks() {
			if !board.ContainsAll(r.board, startBoard) {
				panic(fmt.Errorf(
					"algo %s removed values from the board: before %v, after %v",
					algo, startBoard, r.board))
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
			return status
		}
		// with unknown status all algos must retain the board as is
		if board.GetIntegrityChecks() {
			if !board.Equivalent(r.board, startBoard) {
				panic(fmt.Errorf(
					"algo %s changed the board with unknown status:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, r.board))
			}
		}
	}
	return StatusUnknown
}

func (r *runner) Action() Action {
	return r.action
}
func (r *runner) Board() board.Board {
	return r.board
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

func (r *runner) MergeSteps(steps *StepStats) {
	r.result.Steps.Merge(steps)
}

func (r *runner) recursiveRun(ctx context.Context, b board.Board) *Result {
	nested := r.nestedCache
	if nested == nil {
		nested = &runner{
			action:                r.action,
			board:                 b,
			currentRecursionDepth: r.currentRecursionDepth + 1,
			maxRecursionDepth:     r.maxRecursionDepth,

			// result includes ref to the shared solutions
			result: Result{
				Status:    StatusUnknown,
				Solutions: r.result.Solutions,
			},
			algorithms: r.algorithms,
		}
		r.nestedCache = nested

	} else {
		// we just need to set the board and reset the result's status
		nested.board = b
		nested.currentRecursionDepth = r.currentRecursionDepth + 1
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

	if board.GetIntegrityChecks() {
		// capture the original board for integrity checks to make sure algos do not corrupt
		// the board and their solutions solve the input
		nested.inputBoard = b.Clone(board.Immutable)
	}

	return nested.run(ctx)
}
