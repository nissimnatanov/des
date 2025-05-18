package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"gotest.tools/v3/assert"
)

func TestBitSet81(t *testing.T) {
	var bs indexes.BitSet81

	assert.Assert(t, !bs.Get(0))
	assert.Assert(t, !bs.Get(55))
	assert.Assert(t, !bs.Get(80))

	bs = indexes.MaxBitSet81
	for i := range indexes.BoardSize {
		assert.Assert(t, bs.Get(i))
	}

	bs.Set(80, false)
	assert.Assert(t, !bs.Get(80))
	assert.Assert(t, bs.Get(79))

	bs.Set(80, true)
	assert.Assert(t, bs == indexes.MaxBitSet81)

	bs = indexes.MinBitSet81
	assert.Assert(t, !bs.Get(0))
	assert.Assert(t, !bs.Get(55))
	assert.Assert(t, !bs.Get(80))
}

// BenchmarkBitSet81 tests how the bitSet81BitsPerUnit impacts on
// overall performance of the BitSet81 struct
//
// * 8 (byte): 		0.0004034 ns/op
// * 16 (uint16): 	0.0003733 ns/op (>8MB of static cache memory)
// * 20 (uint32): 	0.0003835 ns/op (>185MB of static cache memory)
// * 15 (uint16):   0.0004229 ns/op (>4MB of static cache memory)
// * 11 (uint16): 	0.0005070 ns/op (>0.2MB of static cache memory)
func BenchmarkBitSet81(b *testing.B) {
	samples := []indexes.BitSet81{
		indexes.MinBitSet81,
		indexes.MaxBitSet81,
	}

	for range 1000 {
		for _, sample := range samples {
			sample.Get(80)
			for i := range sample.Indexes {
				sample.Get(i)
			}
			for i := range sample.Indexes {
				sample.Set(i, sample.Get(i))
			}
			sample.Intersect(indexes.MaxBitSet81)
			sample.Complement().Intersect(indexes.MinBitSet81)
			for range 10 {
				_ = sample.First()
			}
		}
	}
}
