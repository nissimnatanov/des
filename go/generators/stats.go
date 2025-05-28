package generators

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// SolutionStats holds a snapshot of the solution generation statistics.
type SolutionStats struct {
	Count   int64
	Retries int64
	Elapsed time.Duration
}

func (s SolutionStats) AverageElapsed() time.Duration {
	if s.Count == 0 {
		return 0
	}
	return s.Elapsed / time.Duration(s.Count)
}

func (s SolutionStats) AverageRetries() float64 {
	if s.Count == 0 {
		return 0
	}
	return math.Round(float64(s.Retries)*10/float64(s.Count)) / 10
}

func (s SolutionStats) String() string {
	return fmt.Sprintf("Solutions: %d, Average Elapsed: %s, Average Retries: %v",
		s.Count,
		s.AverageElapsed(),
		s.AverageRetries())
}

var Stats stats

type stats struct {
	// lock is shared for now, we can split it later if needed
	rw sync.RWMutex

	solution SolutionStats
}

func (s *stats) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution.Count = 0
	s.solution.Retries = 0
	s.solution.Elapsed = 0
}

func (s *stats) Solution() SolutionStats {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.solution
}

func (s *stats) reportOneSolution(elapsed time.Duration, retries int64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution.Count++
	s.solution.Elapsed += elapsed
	s.solution.Retries += retries
}
