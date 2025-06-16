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
// the recursion is used to solve the board. For safety, set it to the
// MaxFreeCellsForValidBoard, but technically speaking it would never come
// close to that number since the recursion algorithm uses 'the only choice'
// algorithm first and that would detect values.
const maxRecursionDepthLimit = MaxFreeCellsForValidBoard

// Each Solver can be used on a single thread only.
type Solver struct {
	solveAlgorithms     []Algorithm
	proveAlgorithms     []Algorithm
	fastSolveAlgorithms []Algorithm
}

func New() *Solver {
	return &Solver{
		solveAlgorithms:     getAlgorithms(ActionSolve),
		proveAlgorithms:     getAlgorithms(ActionProve),
		fastSolveAlgorithms: getAlgorithms(ActionSolveFast),
	}
}

type options struct {
	actionSet bool
	action    Action
	withSteps bool
}

type Option interface {
	applySolverOptions(*options)
}

func (s *Solver) Run(ctx context.Context, b *boards.Game, os ...Option) *Result {
	if b == nil {
		panic("solver.Run called with nil board")
	}
	var opts options
	for _, o := range os {
		o.applySolverOptions(&opts)
	}
	opts.withSteps = opts.withSteps || boards.GetIntegrityChecks()
	action := opts.action
	// prevent input from being modified
	b = b.Clone(boards.Immutable)
	result := &Result{
		Action: action,
		Input:  b,
	}

	start := time.Now()
	defer func() {
		result.Elapsed = time.Since(start)
	}()

	if action == ActionSolve {
		// ActionSolve uses layered recursion to calculate the level,
		// if the board has more than one solution Solve can take
		// really long time to run. Hence, we do not allow level calculation
		// on an unproven board.
		proofRes := s.Run(ctx, b, ActionProve)
		if proofRes.Status != StatusSucceeded {
			proofRes.Action = action
			return proofRes
		}
		result.addNonLevelSteps(proofRes.Steps)
	}

	runRes := s.run(ctx, b, action, opts.withSteps)
	result.completeWith(runRes)
	logResult(result)
	return result
}

func (s *Solver) run(ctx context.Context, b *boards.Game, action Action, withSteps bool) *runResult {
	var algorithms []Algorithm
	switch action {
	case ActionSolve:
		algorithms = s.solveAlgorithms
	case ActionSolveFast:
		algorithms = s.fastSolveAlgorithms
	case ActionProve:
		algorithms = s.proveAlgorithms
	default:
		panic(fmt.Sprintf("unknown action: %s", action))
	}

	// before we modify the input board, we must clone it into Play mode
	play := b.Clone(boards.Play)

	r := &runner{
		play:                  play,
		algorithms:            algorithms,
		currentRecursionDepth: 0,
		// without recursion it is virtually impossible to solve many boards,
		// zero is not a valid value
		// note: recursion with this package is almost 'allocation-free', and it is fast
		maxRecursionDepth: maxRecursionDepthLimit,
		result:            &runResult{},
		action:            action,
		input:             b,
	}
	if withSteps {
		r.result.Steps = Steps{}
	}

	if !b.IsValid() {
		return r.result.complete(StatusNoSolution)
	}
	if b.IsSolved() {
		// fast track the result generation for board that is already solved (e.g. Solution)
		r.result.Solutions = r.result.Solutions.Append(boards.NewSolution(b))
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

	if len(r.result.Solutions) > 1 && r.result.Status == StatusSucceeded {
		// if we got two solutions we must guarantee a different status
		r.result.complete(StatusMoreThanOneSolution)
	}

	// if run panics, we want to return both the error and the partial result
	return r.result
}
