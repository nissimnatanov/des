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
	Steps Steps `json:"steps"`
}

func (r *Result) complete(status Status) *Result {
	if r.Status != StatusUnknown {
		panic("result already completed")
	}
	r.Status = status
	return r
}

func (r *Result) completeErr(err error) *Result {
	if r.Status != StatusUnknown {
		panic("result already completed")
	}
	r.Status = StatusError
	r.Error = err
	return r
}

func (r *Result) Merge(other *Result) {
	r.Count += other.Count
	r.Complexity += other.Complexity
	r.Level = LevelFromComplexity(r.Complexity)
	r.Steps.Merge(other.Steps)
}

func (r *Result) AddStep(step Step, complexity StepComplexity, count int) {
	switch {
	case count <= 0:
		panic("count must be > 0")
	case complexity <= 0:
		panic("complexity must be > 0")
	case step == "":
		panic("step must not be empty")
	}

	r.Count += int64(count)
	r.Complexity += complexity * StepComplexity(count)
	r.Level = LevelFromComplexity(r.Complexity)
	r.Steps.Add(step, complexity, count)
}

func (r *Result) reset() {
	r.Count = 0
	r.Complexity = 0
	r.Level = 0
	r.Steps = Steps{}
}
