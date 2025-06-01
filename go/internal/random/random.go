package random

import (
	"math/rand"
	"time"
)

func New() *Random {
	return WithSeed(time.Now().UnixNano())
}

func WithSeed(seed int64) *Random {
	r := &Random{
		r:    rand.New(rand.NewSource(seed)),
		seed: seed,
	}
	return r
}

type Random struct {
	r    *rand.Rand
	seed int64
}

func (r *Random) Seed() int64 {
	return r.seed
}

func (r *Random) Intn(n int) int {
	return r.r.Intn(n)
}

func (r *Random) NextInClosedRange(min, max int) int {
	switch {
	case min < max:
		return r.Intn(max-min+1) + min
	case min == max:
		return min
	default:
		panic("min must be less than or equal to max")
	}
}

func (r *Random) PercentProbability(p int) bool {
	return r.Intn(100) < p
}

func Pick[S ~[]T, T any](r *Random, slice S) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	index := r.Intn(len(slice))
	return slice[index], true
}

func Shuffle[S ~[]T, T any](r *Random, s S) {
	if len(s) == 0 {
		return
	}

	r.r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
