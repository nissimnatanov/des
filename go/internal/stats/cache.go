package stats

import "fmt"

type Cache struct {
	HitCount  int64
	MissCount int64
	SetCount  int64

	// for Solve with unknown results only
	UnknownHitCount int64
	UnknownSetCount int64
}

func (s Cache) String() string {
	hitPercent := 0.0
	unknownHitPercent := 0.0
	total := s.HitCount + s.UnknownHitCount + s.MissCount
	if total > 0 {
		hitPercent = float64(s.HitCount) / float64(total) * 100.0
		unknownHitPercent = float64(s.UnknownHitCount) / float64(total) * 100.0
	}
	return fmt.Sprintf("Solver Cache: hits=%d (%.2f%%, unknown=%.2f%%), misses=%d, sets=%d (unknown=%d)",
		s.HitCount, hitPercent, unknownHitPercent, s.MissCount, s.SetCount, s.UnknownSetCount)
}

func (s *Cache) MergeAndDrain(other Cache) {
	s.HitCount += other.HitCount
	s.MissCount += other.MissCount
	s.SetCount += other.SetCount
	s.UnknownHitCount += other.UnknownHitCount
	s.UnknownSetCount += other.UnknownSetCount
	other.HitCount = 0
	other.MissCount = 0
	other.SetCount = 0
	other.UnknownHitCount = 0
	other.UnknownSetCount = 0
}
