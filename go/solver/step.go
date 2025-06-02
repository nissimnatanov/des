package solver

import (
	"encoding/json"

	"github.com/nissimnatanov/des/go/boards"
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

type StepStats struct {
	Count      int64          `json:"count"`
	Complexity StepComplexity `json:"complexity"`
	Level      Level          `json:"level"`

	// Steps are captured only if board's integrity checks are enabled
	Steps map[Step]map[StepComplexity]int `json:"steps,omitempty"`
}

func (s *StepStats) AddStep(step Step, complexity StepComplexity, count int) {
	switch {
	case count <= 0:
		panic("count must be > 0")
	case complexity <= 0:
		panic("complexity must be > 0")
	case step == "":
		panic("step must not be empty")
	}

	s.Count += int64(count)
	s.Complexity += complexity * StepComplexity(count)
	s.Level = LevelFromComplexity(s.Complexity)
	if !boards.GetIntegrityChecks() {
		return
	}
	if s.Steps == nil {
		s.Steps = map[Step]map[StepComplexity]int{}
	}
	if _, ok := s.Steps[step]; !ok {
		s.Steps[step] = map[StepComplexity]int{}
	}
	if _, ok := s.Steps[step][complexity]; !ok {
		s.Steps[step][complexity] = 0
	}
	s.Steps[step][complexity]++
}

func (s *StepStats) String() string {
	str, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(str)
}

func (s *StepStats) Merge(other *StepStats) {
	if other == nil {
		return
	}
	s.Count += other.Count
	s.Complexity += other.Complexity
	s.Level = LevelFromComplexity(s.Complexity)
	if !boards.GetIntegrityChecks() || len(other.Steps) == 0 {
		return
	}
	if s.Steps == nil {
		s.Steps = map[Step]map[StepComplexity]int{}
	}
	for step, complexityMap := range other.Steps {
		if _, ok := s.Steps[step]; !ok {
			s.Steps[step] = map[StepComplexity]int{}
		}
		for complexity, count := range complexityMap {
			s.Steps[step][complexity] += count
		}
	}
}

func (s *StepStats) reset() {
	s.Count = 0
	s.Complexity = 0
	s.Level = Level(0)
	if s.Steps != nil {
		clear(s.Steps)
	}
}
