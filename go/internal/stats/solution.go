package stats

import (
	"fmt"
	"math"
	"time"
)

// Solution holds a snapshot of the solution generation statistics.
type Solution struct {
	Count   int64
	Retries int64
	Elapsed time.Duration
}

func (s Solution) AverageElapsed() time.Duration {
	if s.Count == 0 {
		return 0
	}
	return s.Elapsed / time.Duration(s.Count)
}

func (s Solution) AverageRetries() float64 {
	if s.Count == 0 {
		return 0
	}
	return math.Round(float64(s.Retries)*10/float64(s.Count)) / 10
}

func (s Solution) String() string {
	return fmt.Sprintf("Solutions: %d, ~Elapsed: %s, ~Retries: %.3f",
		s.Count,
		s.AverageElapsed(),
		s.AverageRetries())
}

func (s *Solution) reportOne(elapsed time.Duration, retries int64) {
	s.Count++
	s.Elapsed += elapsed
	s.Retries += retries
}
