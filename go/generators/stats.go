package generators

import (
	"sync"
	"time"
)

var Stats stats

type stats struct {
	// lock is shared for now, we can split it later if needed
	rw sync.RWMutex

	solution SolutionStats
	game     GameStats
}

func (s *stats) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution = SolutionStats{}
	s.game = GameStats{}
}

func (s *stats) Solution() SolutionStats {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.solution
}

func (s *stats) Game() GameStats {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.game.clone()
}

func (s *stats) reportOneSolution(elapsed time.Duration, retries int64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution.reportOne(elapsed, retries)
}

func (s *stats) reportOneGeneration(elapsed time.Duration, retries int64, stageStats GamePerStageStats) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.game.reportOne(elapsed, retries, stageStats)
}
