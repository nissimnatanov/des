package solver

import (
	"time"

	"github.com/nissimnatanov/des/go/boards"
)

type Result struct {
	// Input is the initial input board, as provided by the caller.
	Input      *boards.Game   `json:"board,omitempty"`
	Action     Action         `json:"action,omitempty"`
	Status     Status         `json:"status"`
	Error      error          `json:"error,omitempty"`
	Solutions  Solutions      `json:"-"`
	Elapsed    time.Duration  `json:"elapsed"`
	Count      int64          `json:"count"`
	Complexity StepComplexity `json:"complexity"`
	Level      Level          `json:"level"`
	// Steps are the steps that led to the current level.
	Steps Steps `json:"level_steps"`

	// AllSteps includes all all the steps performed by the solver, including those that
	// are not directly related to the current level, but had to be performed (like Proof before
	// the actual solving, or layered recursion).
	AllSteps Steps `json:"all_steps"`
}

func (r *Result) addNonLevelSteps(steps Steps) {
	r.AllSteps.Merge(steps)
}

func (r *Result) completeWith(runRes *runResult) *Result {
	if r.Status != StatusUnknown {
		panic("result already completed")
	}
	if r.Count != 0 || r.Complexity != 0 {
		panic("result should not have any partial count nor complexity reported, they are override below")
	}
	if runRes == nil || runRes.Status == StatusUnknown {
		panic("solver returned nil or unknown status")
	}

	r.Status = runRes.Status
	r.Error = runRes.Error
	r.Count = runRes.Count
	r.Complexity = runRes.Complexity
	r.Solutions = r.Solutions.With(runRes.Solutions)
	if r.Status == StatusSucceeded {
		r.Level = LevelFromComplexity(r.Complexity)
	}
	r.Steps.Merge(runRes.Steps)
	// Note: AllSteps might already have some results, we should always merge with them
	r.AllSteps.Merge(runRes.Steps)
	return r
}
