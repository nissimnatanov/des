package solver

import "fmt"

type Status int

const (
	StatusUnknown Status = iota
	StatusSucceeded
	StatusNoSolution
	StatusMoreThanOneSolution
	StatusLessThan17
	StatusTwoOrMoreValuesMissing
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "Unknown"
	case StatusSucceeded:
		return "Succeeded"
	case StatusNoSolution:
		return "No solution"
	case StatusMoreThanOneSolution:
		return "Two or more solutions"
	case StatusLessThan17:
		return "Less than 17"
	case StatusError:
		return "Error or panic"
	case StatusTwoOrMoreValuesMissing:
		return "Two or more values missing"
	default:
		return fmt.Sprintf("WRONG SudokuSolverStatus: %d", s)
	}
}
