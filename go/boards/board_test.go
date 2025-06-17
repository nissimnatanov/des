package boards_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/boards/values"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestEmptyBoard(t *testing.T) {
	b := boards.New()

	assert.Equal(t, boards.Edit, b.Mode())
	assert.Assert(t, b.IsValid())
	assert.Assert(t, b.IsValidCell(2))
	assert.Assert(t, b.IsValidCell(80))
	assert.Assert(t, !b.IsSolved())
	assert.Assert(t, b.IsEmpty(33))
	assert.Equal(t, values.Value(0), b.Get(0))
	assert.Equal(t, values.Value(0), b.Get(80))
	assert.Equal(t, 81, b.FreeCellCount())

	assert.Equal(t, values.FullSet, b.AllowedValues(0))
	assert.Equal(t, values.FullSet, b.AllowedValues(80))

	assert.Assert(t, boards.ContainsAll(b, b))
	assert.Assert(t, boards.ContainsReadOnly(b, b))
	assert.Assert(t, boards.Equivalent(b, b))
	assert.Assert(t, boards.EquivalentReadOnly(b, b))

	assert.Equal(t, "ZZZC", boards.Serialize(b))

	// negative assertions
	assert.Assert(t, cmp.Panics(func() { b.IsValidCell(81) }))
	assert.Assert(t, cmp.Panics(func() { b.Get(81) }))
	assert.Assert(t, cmp.Panics(func() { b.Get(-2) }))
	assert.Assert(t, cmp.Panics(func() { b.AllowedValues(-1) }))
	assert.Assert(t, cmp.Panics(func() { b.AllowedValues(81) }))
	assert.Assert(t, cmp.Panics(func() { b.DisallowValue(81, 7) }))
	assert.Assert(t, cmp.Panics(func() { b.DisallowValues(-1, values.Value(8).AsSet()) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(-1, values.Value(8)) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(81, values.Value(5)) }))
	assert.Assert(t, cmp.Panics(func() { b.Set(81, values.Value(4)) }))

	// special cases
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(0, 0) }))
	assert.Assert(t, cmp.Panics(func() { b.DisallowValue(0, 0) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(0, 10) }))
}

func newSampleBoard() *boards.Game {
	b := boards.New()

	col0 := 0
	col1 := 3
	col2 := 6

	for v := values.Value(1); v <= 9; v++ {
		b.SetReadOnly(col0, v)
		b.Set(9+col1, v)
		b.Set(18+col2, v)
		col0 = (col0 + 1) % boards.SequenceSize
		col1 = (col1 + 1) % boards.SequenceSize
		col2 = (col2 + 1) % boards.SequenceSize
	}
	return b
}

func assertSampleBoard(t *testing.T, b *boards.Game) {
	assert.Assert(t, b.IsValid())
	assert.Assert(t, b.IsValidCell(4))
	assert.Assert(t, b.IsValidCell(40))
	assert.Assert(t, boards.Equivalent(b, newSampleBoard()), "Boards are not equal", b)
}

func TestValidBoard(t *testing.T) {
	b := newSampleBoard()
	// cspell:disable-next-line
	bs := boards.Format(b, "vt")
	t.Log("\n" + bs)
	t.Fail()

	assertSampleBoard(t, b)

	b.Set(40, values.Value(5))
	assert.Assert(t, !b.IsValid())

	assert.Assert(t, !b.IsValidCell(40))
	assert.Assert(t, b.IsValidCell(4)) // single read-only

	b.Set(40, 0)
	assertSampleBoard(t, b)
}

func TestClone(t *testing.T) {
	b := newSampleBoard()

	play := b.Clone(boards.Play)
	assertSampleBoard(t, play)
	assert.Equal(t, boards.Play, play.Mode())

	edit := b.Clone(boards.Edit)
	assertSampleBoard(t, edit)
	assert.Equal(t, boards.Edit, edit.Mode())
}

func TestPlay(t *testing.T) {
	b := newSampleBoard()
	play := b.Clone(boards.Play)
	assertSampleBoard(t, play)
	assert.Equal(t, boards.Play, play.Mode())

	assert.Assert(t, cmp.Panics(func() { play.SetReadOnly(55, values.Value(6)) }))
	assert.Equal(t, values.Value(0), play.Get(55))
	assertSampleBoard(t, play)

	play.Set(55, values.Value(6))
	assert.Equal(t, values.Value(6), play.Get(55))
	assert.Assert(t, play.IsValid())
	assert.Assert(t, boards.ContainsReadOnly(play, b))
	assert.Assert(t, boards.ContainsAll(play, b))

	play.Set(40, values.Value(5))
	assert.Equal(t, values.Value(5), play.Get(40))
	assert.Assert(t, !play.IsValid())
	assert.Assert(t, boards.ContainsReadOnly(play, b))
	assert.Assert(t, boards.ContainsAll(play, b))

	play.Set(40, values.Value(9))
	assert.Assert(t, play.IsValid())
	assert.Assert(t, boards.ContainsReadOnly(play, b))
	assert.Assert(t, boards.ContainsAll(play, b))

	play.Restart()
	assert.Equal(t, values.Value(0), play.Get(55))
	assert.Equal(t, values.Value(0), play.Get(40))
	// only the first row is read-only
	assert.Assert(t, boards.ContainsReadOnly(b, play))
}
