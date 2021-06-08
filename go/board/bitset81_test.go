package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitSet81(t *testing.T) {
	var bs BitSet81
	assert := assert.New(t)

	assert.False(bs.Get(0))
	assert.False(bs.Get(55))
	assert.False(bs.Get(80))
	assert.False(bs.AllSet())

	bs.SetAll(true)
	assert.True(bs.Get(0))
	assert.True(bs.Get(44))
	assert.True(bs.Get(80))
	assert.True(bs.AllSet())

	bs.Set(80, false)
	assert.False(bs.AllSet())
	assert.False(bs.Get(80))
	assert.True(bs.Get(79))

	bs.Set(80, true)
	assert.True(bs.AllSet())

	bs.Set(33, false)
	assert.False(bs.AllSet())

	bs.Reset()
	assert.False(bs.Get(0))
	assert.False(bs.Get(55))
	assert.False(bs.Get(80))
}
