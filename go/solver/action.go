package solver

import "fmt"

type Action int

const (
	ActionProve Action = iota
	ActionSolveFast
	// ActionSolve always runs Prove first to make sure solver does not
	// run too long on a board that has too many solutions, entering very
	// deep recursion.
	ActionSolve
)

func (a Action) String() string {
	switch a {
	case ActionProve:
		return "Prove"
	case ActionSolveFast:
		return "FastSolve"
	case ActionSolve:
		return "Solve"
	// case ActionHint:
	// 	return "Hint"
	default:
		return fmt.Sprintf("WRONG SudokuSolverAction: %d", a)
	}
}

// LevelRequested returns true if the action requires an accurate level, e.g. solver and
// its algorithms should mimic human solving techniques as much as possible
func (a Action) LevelRequested() bool {
	return a == ActionSolve
}

// ProofRequested returns true if the action requires a proof of the solution, e.g. solver and
// its algorithms should prove the solution is unique
func (a Action) ProofRequested() bool {
	return a == ActionProve
}

func (a Action) applySolverOptions(o *options) {
	if o.actionSet {
		panic("do not provide action twice")
	}
	o.action = a
	o.actionSet = true
}
