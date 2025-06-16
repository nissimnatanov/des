package solver

import (
	"encoding/json"
)

// Steps are reported by each algorithm
type Step string

type StepComplexity int64

const (
	StepComplexityEasy       StepComplexity = 1       // single in square
	StepComplexityMedium     StepComplexity = 5       // single in row/column
	StepComplexityHard       StepComplexity = 20      // identify pairs
	StepComplexityHarder     StepComplexity = 50      // identify triplets
	StepComplexityRecursion1 StepComplexity = 100     // single recursion or any recursion if level is not needed
	StepComplexityRecursion2 StepComplexity = 10000   // second-level is rare, it usually leads to a Black Hole level
	StepComplexityRecursion3 StepComplexity = 100000  // third-level recursion - not reached yet with layered recursion
	StepComplexityRecursion4 StepComplexity = 1000000 // fourth-level recursion and beyond
)

// Steps are captured only if board's integrity checks are enabled
type Steps map[Step]map[StepComplexity]int

func (s Steps) Add(step Step, complexity StepComplexity, count int) {
	// Steps is allocated only if requested
	if s == nil {
		return
	}

	switch {
	case count <= 0:
		panic("count must be > 0")
	case complexity <= 0:
		panic("complexity must be > 0")
	case step == "":
		panic("step must not be empty")
	}

	if _, ok := s[step]; !ok {
		s[step] = map[StepComplexity]int{}
	}
	if _, ok := s[step][complexity]; !ok {
		s[step][complexity] = 0
	}
	s[step][complexity]++
}

func (s Steps) String() string {
	if s == nil {
		return "{}"
	}
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(str)
}

func (s Steps) Merge(other Steps) {
	if s == nil || len(other) == 0 {
		return
	}
	for step, complexityMap := range other {
		if _, ok := s[step]; !ok {
			s[step] = map[StepComplexity]int{}
		}
		for complexity, count := range complexityMap {
			s[step][complexity] += count
		}
	}
}

// WithSteps can be added to the Solver as an option to capture the step stats
func WithSteps(opts *options) {
	opts.withSteps = true
}
