package solver

import (
	"maps"
	"slices"
)

type runResult struct {
	Status     Status         `json:"status"`
	Count      int64          `json:"count"`
	Complexity StepComplexity `json:"complexity"`
	Error      error          `json:"error,omitempty"`
	Steps      Steps          `json:"steps"`
	Solutions  Solutions      `json:"-"`
}

func (rr *runResult) complete(status Status) *runResult {
	rr.Status = status
	return rr
}

func (rr *runResult) completeErr(err error) *runResult {
	rr.Status = StatusError
	rr.Error = err
	return rr
}

func (rr *runResult) clone() *runResult {
	c := *rr
	if len(rr.Steps) != 0 {
		c.Steps = make(map[Step]map[StepComplexity]int, len(rr.Steps))
		for step, complexityMap := range rr.Steps {
			// values of the nested maps are values, shallow clone is good
			c.Steps[step] = maps.Clone(complexityMap)
		}
	}
	if len(c.Solutions) != 0 {
		c.Solutions = slices.Clone(c.Solutions)
	}
	return &c
}

func (rr *runResult) addStep(step Step, complexity StepComplexity, count int) {
	switch {
	case count <= 0:
		panic("count must be > 0")
	case complexity <= 0:
		panic("complexity must be > 0")
	case step == "":
		panic("step must not be empty")
	}

	rr.Count += int64(count)
	rr.Complexity += complexity * StepComplexity(count)
	rr.Steps.Add(step, complexity, count)
}

// mergeStatsOnly is called multiple times, with various test values tried by the
// trialAndError algo, we only merge the basic stats/steps here, but not the status,
// relying on the trialAndError to report the final status based on its findings
func (rr *runResult) mergeStatsOnly(other *runResult) {
	rr.Count += other.Count
	rr.Complexity += other.Complexity
	rr.Steps.Merge(other.Steps)
	rr.Solutions = rr.Solutions.With(other.Solutions)
}
