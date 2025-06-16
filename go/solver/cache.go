package solver

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/nissimnatanov/des/go/boards"
)

type CacheKey string

type CacheValue struct {
	result *runResult
	board  *boards.Game
}

func (cv CacheValue) IsPresent() bool {
	return cv.result != nil
}

func (cv CacheValue) clone() CacheValue {
	if !cv.IsPresent() {
		panic("cannot clone an empty cache value")
	}
	clone := CacheValue{
		result: cv.result.clone(),
	}
	if cv.board != nil {
		// no-op if already immutable
		clone.board = cv.board.Clone(boards.Immutable)
	}
	return clone
}

type cacheValue struct {
	// solveCVs maps recursion depth tried to its Solve results, so we can return the exact match
	// it is only used for Solve which uses layered recursion and needs accurate leveling
	solveCVs [4]CacheValue

	proveCV CacheValue
}

var NoCache = CacheKey("")

func NewCache() *Cache {
	return &Cache{
		m: make(map[CacheKey]*cacheValue),
	}
}

type Cache struct {
	m map[CacheKey]*cacheValue

	hitCount  atomic.Int64
	missCount atomic.Int64
	setCount  atomic.Int64

	// hits for unknown results, only for Solve
	unknownHitCount atomic.Int64
}

type CacheStats struct {
	HitCount  int64
	MissCount int64
	SetCount  int64

	// for Solve with unknown results only
	UnknownHitCount int64
}

func (s CacheStats) String() string {
	hitPercent := 0.0
	unknownHitPercent := 0.0
	total := s.HitCount + s.UnknownHitCount + s.MissCount
	if total > 0 {
		hitPercent = float64(s.HitCount) / float64(total) * 100.0
		unknownHitPercent = float64(s.UnknownHitCount) / float64(total) * 100.0
	}
	return fmt.Sprintf("Solver Cache: hits=%d (%.1f%%), unknown hits=%d (%.1f%%), misses=%d, sets=%d",
		s.HitCount, hitPercent, s.UnknownHitCount, unknownHitPercent, s.MissCount, s.SetCount)
}

func (s *CacheStats) MergeAndDrain(other CacheStats) {
	s.HitCount += other.HitCount
	s.MissCount += other.MissCount
	s.SetCount += other.SetCount
	s.UnknownHitCount += other.UnknownHitCount
	other.HitCount = 0
	other.MissCount = 0
	other.SetCount = 0
	other.UnknownHitCount = 0
}

func (c *Cache) makeKey(b *boards.Game) CacheKey {
	if c == nil {
		return NoCache
	}
	sb := strings.Builder{}
	for i, v := range b.AllValues {
		sb.WriteRune('0' + rune(v))
		if v != 0 {
			continue
		}
		disallowedByUser := b.DisallowedByUser(i)
		if !disallowedByUser.IsEmpty() {
			sb.WriteRune('[')
			for _, d := range disallowedByUser.Values() {
				sb.WriteRune('0' + rune(d))
			}
			sb.WriteRune(']')
		}
	}
	return CacheKey(sb.String())
}

func (c *Cache) get(b *boards.Game, action Action, maxRecursionDepthRemained int) (CacheValue, CacheKey) {
	if c == nil {
		return CacheValue{}, NoCache
	}
	key := c.makeKey(b)
	r := c.m[key]
	if r == nil {
		c.missCount.Add(1)
		return CacheValue{}, key
	}

	// we have cache for this board, let's see what the action was
	switch action {
	case ActionSolve:
		// keep going, Solve needs special handling for layered results
	case ActionProve, ActionSolveFast:
		// do not try re-using Solve results for Prove or SolveFast, reusing Solve results may
		// actually be much slower since solver uses layered recursion and progresses very slow
		// to emulate human solving
		if r.proveCV.IsPresent() {
			c.hitCount.Add(1)
			return r.proveCV.clone(), key
		}
		// do not try using unknown results for proving since they have limited recursion depth
		return CacheValue{}, key
	default:
		panic(fmt.Sprintf("unknown action in cache: %s", action))
	}

	// action is Solve and solveRes is not known yet
	if maxRecursionDepthRemained > len(r.solveCVs) || !r.solveCVs[maxRecursionDepthRemained].IsPresent() {
		c.missCount.Add(1)
		return CacheValue{}, key
	}
	cv := r.solveCVs[maxRecursionDepthRemained]
	if cv.result.Status == StatusUnknown {
		c.unknownHitCount.Add(1)
	} else {
		c.hitCount.Add(1)
	}
	return cv.clone(), key
}

func (c *Cache) set(key CacheKey, action Action, cv CacheValue, maxRecursionDepthTried int) {
	if c == nil {
		return
	}
	if key == NoCache {
		panic("cannot set cache with empty key if cache is enabled")
	}
	if action == ActionSolveFast {
		// do not care about caching SolveFast, it is not used for generation
		return
	}
	c.setCount.Add(1)
	r := c.m[key]
	if r == nil {
		r = &cacheValue{}
		c.m[key] = r
	}
	// clone the cache value as immutable for caching
	cv = cv.clone()

	if action == ActionProve {
		// prove does not use layered recursion and its result is always deterministic
		r.proveCV = cv
		return
	}

	if action != ActionSolve {
		panic(fmt.Sprintf("unknown action in cache: %s", action))
	}

	// result is unknown, let's capture the max recursion depth tried so far
	if maxRecursionDepthTried > len(r.solveCVs) {
		// we only support up to 4 recursion depths since we prove boards first
		return
	}
	// we do not need the board for unknown results, it is always the same as the input
	r.solveCVs[maxRecursionDepthTried] = cv
}

func (c *Cache) Stats() CacheStats {
	if c == nil {
		return CacheStats{}
	}
	return CacheStats{
		HitCount:        c.hitCount.Load(),
		MissCount:       c.missCount.Load(),
		SetCount:        c.setCount.Load(),
		UnknownHitCount: c.unknownHitCount.Load(),
	}
}

func (c *Cache) ResetStats() {
	if c == nil {
		return
	}
	c.hitCount.Store(0)
	c.missCount.Store(0)
	c.setCount.Store(0)
	c.unknownHitCount.Store(0)
}

func (c *Cache) applySolverOptions(opts *options) {
	if opts.cache != nil {
		panic("cannot set cache options twice")
	}
	opts.cache = c
}
