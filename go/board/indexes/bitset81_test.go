package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board/indexes"
	"gotest.tools/v3/assert"
)

func TestBitSet81(t *testing.T) {
	var bs indexes.BitSet81

	assert.Assert(t, !bs.Get(0))
	assert.Assert(t, !bs.Get(55))
	assert.Assert(t, !bs.Get(80))
	assert.Assert(t, !bs.AllSet())

	bs.SetAll(true)
	assert.Assert(t, bs.Get(0))
	assert.Assert(t, bs.Get(44))
	assert.Assert(t, bs.Get(80))
	assert.Assert(t, bs.AllSet())

	bs.Set(80, false)
	assert.Assert(t, !bs.AllSet())
	assert.Assert(t, !bs.Get(80))
	assert.Assert(t, bs.Get(79))

	bs.Set(80, true)
	assert.Assert(t, bs.AllSet())

	bs.Set(33, false)
	assert.Assert(t, !bs.AllSet())

	bs.Reset()
	assert.Assert(t, !bs.Get(0))
	assert.Assert(t, !bs.Get(55))
	assert.Assert(t, !bs.Get(80))
}
