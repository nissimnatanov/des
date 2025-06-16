package internal

import (
	"sync"
	"time"

	"github.com/nissimnatanov/des/go/solver"
)

var Stats stats

type stats struct {
	// lock is shared for now, we can split it later if needed
	rw sync.RWMutex

	solution SolutionStats
	game     GameStats
	cache    solver.CacheStats
}

func (s *stats) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution = SolutionStats{}
	s.game = GameStats{}
	s.cache = solver.CacheStats{}
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

func (s *stats) Cache() solver.CacheStats {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.cache
}

func (s *stats) ReportOneSolution(elapsed time.Duration, retries int64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution.reportOne(elapsed, retries)
}

func (s *stats) ReportGeneration(
	count int, elapsed time.Duration, retries int64,
	stageStats GamePerStageStats, cacheStats solver.CacheStats,
) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.game.report(count, elapsed, retries, stageStats)
	s.cache.MergeAndDrain(cacheStats)
}
