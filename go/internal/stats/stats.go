package stats

import (
	"sync"
	"time"
)

var Stats all

type all struct {
	// lock is shared for now, we can split it later if needed
	rw sync.RWMutex

	solution Solution
	game     GameStats
	cache    Cache
}

func (s *all) Reset() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution = Solution{}
	s.game = GameStats{}
	s.cache = Cache{}
}

func (s *all) Solution() Solution {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.solution
}

func (s *all) Game() GameStats {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.game.clone()
}

func (s *all) Cache() Cache {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.cache
}

func (s *all) ReportOneSolution(elapsed time.Duration, retries int64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.solution.reportOne(elapsed, retries)
}

func (s *all) ReportGeneration(
	count int, elapsed time.Duration, retries int64,
	stageStats GameStages, cacheStats Cache,
) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.game.report(count, elapsed, retries, stageStats)
	s.cache.MergeAndDrain(cacheStats)
}
