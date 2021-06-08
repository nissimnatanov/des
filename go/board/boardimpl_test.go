package board

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyBoard(t *testing.T) {
	assert := assert.New(t)
	b := NewBoard()

	assert.Equal(Edit, b.Mode())
	assert.True(b.IsValid())
	assert.True(b.IsValidCell(2))
	assert.True(b.IsValidCell(80))
	assert.False(b.IsSolved())
	assert.True(b.IsEmpty(33))
	assert.Equal(Empty, b.Get(0))
	assert.Equal(Empty, b.Get(80))
	assert.Equal(EmptySet(), b.RowSet(0))
	assert.Equal(EmptySet(), b.ColumnSet(8))
	assert.Equal(EmptySet(), b.SquareSet(2))
	assert.Equal(81, b.FreeCellCount())
	assert.Zero(b.Count(Seven))

	assert.Equal(FullSet(), b.AllowedSet(0))
	assert.Equal(FullSet(), b.AllowedSet(80))

	assert.True(b.ContainsAll(b))
	assert.True(b.ContainsReadOnly(b))
	assert.True(b.IsEquivalent(b))
	assert.True(b.IsEquivalentReadOnly(b))

	// each Z => 27 empty cells (3*27 = 81)
	assert.Equal("ZZZ", Serialize(b))

	// negative assertions
	assert.Panics(func() { b.IsValidCell(81) })
	assert.Panics(func() { b.Get(81) })
	assert.Panics(func() { b.Get(-2) })
	assert.Panics(func() { b.RowSet(-3) })
	assert.Panics(func() { b.ColumnSet(9) })
	assert.Panics(func() { b.SquareSet(111) })
	assert.Panics(func() { b.Count(-1) })
	assert.Panics(func() { b.Count(10) })
	assert.Panics(func() { b.AllowedSet(-1) })
	assert.Panics(func() { b.AllowedSet(81) })
	assert.Panics(func() { b.Disallow(81, Seven) })
	assert.Panics(func() { b.DisallowSet(-1, Eight.AsSet()) })
	assert.Panics(func() { b.SetReadOnly(-1, Eight) })
	assert.Panics(func() { b.SetReadOnly(81, Five) })
	assert.Panics(func() { b.Set(81, Four) })

	// special cases
	assert.Panics(func() { b.SetReadOnly(0, Empty) })
	assert.Panics(func() { b.Disallow(0, Empty) })
	assert.Panics(func() { b.SetReadOnly(0, 10) })
}

func newSampleBoard() Board {
	b := NewBoard()

	col0 := 0
	col1 := 3
	col2 := 6

	for vi := FullSet().Iterator(); vi.Next(); {
		v := vi.Value()
		b.SetReadOnly(col0, v)
		b.Set(9+col1, v)
		b.Set(18+col2, v)
		col0 = (col0 + 1) % SequenceSize
		col1 = (col1 + 1) % SequenceSize
		col2 = (col2 + 1) % SequenceSize
	}
	return b
}

func assertSampleBoard(t *testing.T, b Board) {
	assert := assert.New(t)

	assert.True(b.IsValid())
	assert.True(b.IsValidCell(4))
	assert.True(b.IsValidCell(40))
	if !b.IsEquivalent(newSampleBoard()) {
		assert.Fail("Board is not equivalent", b)
	}
}

func TestValidBoard(t *testing.T) {
	assert := assert.New(t)
	b := newSampleBoard()
	Write(b, NewWriter(func(s string) { t.Log(strings.ReplaceAll(s, "\n", "")) }), "vrcst")
	assertSampleBoard(t, b)

	b.Set(40, Five)
	assert.False(b.IsValid())

	assert.False(b.IsValidCell(40))
	assert.True(b.IsValidCell(4)) // single read-only

	b.Set(40, Empty)
	assertSampleBoard(t, b)
}

func TestClone(t *testing.T) {
	assert := assert.New(t)
	b := newSampleBoard()

	play := b.Clone(Play)
	assertSampleBoard(t, play)
	assert.Equal(Play, play.Mode())

	edit := b.Clone(Edit)
	assertSampleBoard(t, edit)
	assert.Equal(Edit, edit.Mode())
}

func TestPlay(t *testing.T) {
	assert := assert.New(t)
	b := newSampleBoard()
	play := b.Clone(Play)
	assertSampleBoard(t, play)
	assert.Equal(Play, play.Mode())

	assert.Panics(func() { play.SetReadOnly(55, Six) })
	assert.Equal(Empty, play.Get(55))
	assertSampleBoard(t, play)

	play.Set(55, Six)
	assert.Equal(Six, play.Get(55))
	assert.True(play.IsValid())
	assert.True(play.ContainsReadOnly(b))
	assert.True(play.ContainsAll(b))

	play.Set(40, Five)
	assert.Equal(Five, play.Get(40))
	assert.False(play.IsValid())
	assert.True(play.ContainsReadOnly(b))
	assert.True(play.ContainsAll(b))

	play.Set(40, Nine)
	assert.True(play.IsValid())
	assert.True(play.ContainsReadOnly(b))
	assert.True(play.ContainsAll(b))

	play.Restart()
	assert.Equal(Empty, play.Get(55))
	assert.Equal(Empty, play.Get(40))
	// only the first row is read-only
	assert.Equal(FullSet(), play.RowSet(0))
	for row := 1; row < 9; row++ {
		assert.Equal(EmptySet(), play.RowSet(row))
	}
	assert.True(b.ContainsReadOnly(play))
}
