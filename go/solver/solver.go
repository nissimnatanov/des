package solver

import (
	"context"
	"fmt"
	"runtime/debug"

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

func (s *Solver) Run(ctx context.Context, b *boards.Game) *Result {
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
	if boards.GetIntegrityChecks() {
		// capture the original board for integrity checks to make sure algos do not corrupt
		// the board and their solutions solve the input
		r.inputBoard = b.Clone(boards.Immutable)
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
	board                 *boards.Game
	currentRecursionDepth int8
	maxRecursionDepth     int8
	result                Result

	// inputBoard for integrity checks of the solutions
	inputBoard *boards.Game

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
	r.result.Solutions.Add(boards.NewSolution(r.board))
	// we got exactly one solution, and we are done
	if boards.GetIntegrityChecks() {
		for si := range r.result.Solutions.Size() {
			sol := r.result.Solutions.At(si)
			if !boards.ContainsAll(sol, r.inputBoard) {
				panic("solution does not match the board to be solved")
			}
		}
	}
	if r.result.Solutions.Size() > 1 {
		return r.result.complete(StatusMoreThanOneSolution)
	}
	return r.result.complete(StatusSucceeded)
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
			startBoard = r.board.Clone(boards.Immutable)
		}

		freeBefore := r.board.FreeCellCount()
		status := algo.Run(ctx, r)

		if boards.GetIntegrityChecks() {
			if !boards.ContainsAll(r.board, startBoard) {
				panic(fmt.Errorf(
					"algo %s removed values from the board: before %v, after %v",
					algo, startBoard, r.board))
			}
			if !r.board.IsValid() {
				panic(fmt.Errorf(
					"algo %s generated failed board:\nbefore:\n%v\nafter:\n%v",
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
			if status == StatusSucceeded &&
				r.board.FreeCellCount() == freeBefore &&
				!r.action.LevelRequested() &&
				r.board.Hint01() < 0 {
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
			if !boards.Equivalent(r.board, startBoard) {
				panic(fmt.Errorf(
					"algo %s changed the board with unknown status:\nbefore:\n%v\nafter:\n%v",
					algo, startBoard, r.board))
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

func (r *runner) recursiveRun(ctx context.Context, b *boards.Game) *Result {
	nested := r.nestedCache
	if nested == nil {
		nested = &runner{
			action:                r.action,
			currentRecursionDepth: r.currentRecursionDepth + 1,
			maxRecursionDepth:     r.maxRecursionDepth,
			algorithms:            r.algorithms,

			board: b,
			// result includes ref to the shared solutions
			result: Result{
				Status:    StatusUnknown,
				Solutions: r.result.Solutions,
			},
		}
		r.nestedCache = nested
	} else {
		// we just need to set the board and reset the result's status
		nested.board = b
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
		// capture the original board for integrity checks to make sure algos do not corrupt
		// the board and their solutions solve the input
		nested.inputBoard = b.Clone(boards.Immutable)
	}

	return nested.run(ctx)
}
