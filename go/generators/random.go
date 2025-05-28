package generators

import (
	"math/rand"
	"time"
)

func NewRandom() *Random {
	seed := time.Now().UnixNano()
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

/*
func (r *random) percentProbability(p int) bool {
	return r.r.Intn(100) < p
}
*/

func RandPick[S ~[]T, T any](r *Random, slice S) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	index := r.r.Intn(len(slice))
	return slice[index], true
}

func RandShuffle[S ~[]T, T any](r *Random, s S) {
	if len(s) == 0 {
		return
	}

	r.r.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
