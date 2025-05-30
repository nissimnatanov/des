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

	bs.SetTo(80, false)
	assert.Assert(t, !bs.Get(80))
	assert.Assert(t, bs.Get(79))

	bs.SetTo(80, true)
	assert.Assert(t, bs == indexes.MaxBitSet81)

	bs = indexes.MinBitSet81
	assert.Assert(t, !bs.Get(0))
	assert.Assert(t, !bs.Get(55))
	assert.Assert(t, !bs.Get(80))
}
