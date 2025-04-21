package board_test

import (
	"strings"
	"testing"

	"github.com/nissimnatanov/des/go/board"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestEmptyBoard(t *testing.T) {
	b := board.New()

	assert.Equal(t, board.Edit, b.Mode())
	assert.Assert(t, b.IsValid())
	assert.Assert(t, b.IsValidCell(2))
	assert.Assert(t, b.IsValidCell(80))
	assert.Assert(t, !b.IsSolved())
	assert.Assert(t, b.IsEmpty(33))
	assert.Equal(t, board.Value(0), b.Get(0))
	assert.Equal(t, board.Value(0), b.Get(80))
	assert.Equal(t, board.EmptyValueSet(), b.RowSet(0))
	assert.Equal(t, board.EmptyValueSet(), b.ColumnSet(8))
	assert.Equal(t, board.EmptyValueSet(), b.SquareSet(2))
	assert.Equal(t, 81, b.FreeCellCount())
	assert.Equal(t, 0, b.Count(7))

	assert.Equal(t, board.FullValueSet(), b.AllowedSet(0))
	assert.Equal(t, board.FullValueSet(), b.AllowedSet(80))

	assert.Assert(t, board.ContainsAll(b, b))
	assert.Assert(t, board.ContainsReadOnly(b, b))
	assert.Assert(t, board.Equivalent(b, b))
	assert.Assert(t, board.EquivalentReadOnly(b, b))

	// each Z => 27 empty cells (3*27 = 81)
	assert.Equal(t, "ZZZ", board.Serialize(b))

	// negative assertions
	assert.Assert(t, cmp.Panics(func() { b.IsValidCell(81) }))
	assert.Assert(t, cmp.Panics(func() { b.Get(81) }))
	assert.Assert(t, cmp.Panics(func() { b.Get(-2) }))
	assert.Assert(t, cmp.Panics(func() { b.RowSet(-3) }))
	assert.Assert(t, cmp.Panics(func() { b.ColumnSet(9) }))
	assert.Assert(t, cmp.Panics(func() { b.SquareSet(111) }))
	assert.Assert(t, cmp.Panics(func() { b.Count(-1) }))
	assert.Assert(t, cmp.Panics(func() { b.Count(10) }))
	assert.Assert(t, cmp.Panics(func() { b.AllowedSet(-1) }))
	assert.Assert(t, cmp.Panics(func() { b.AllowedSet(81) }))
	assert.Assert(t, cmp.Panics(func() { b.Disallow(81, 7) }))
	assert.Assert(t, cmp.Panics(func() { b.DisallowSet(-1, board.Value(8).AsSet()) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(-1, board.Value(8)) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(81, board.Value(5)) }))
	assert.Assert(t, cmp.Panics(func() { b.Set(81, board.Value(4)) }))

	// special cases
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(0, 0) }))
	assert.Assert(t, cmp.Panics(func() { b.Disallow(0, 0) }))
	assert.Assert(t, cmp.Panics(func() { b.SetReadOnly(0, 10) }))
}

func newSampleBoard() board.Board {
	b := board.New()

	col0 := 0
	col1 := 3
	col2 := 6

	for vi := board.FullValueSet().Iterator(); vi.Next(); {
		v := vi.Value()
		b.SetReadOnly(col0, v)
		b.Set(9+col1, v)
		b.Set(18+col2, v)
		col0 = (col0 + 1) % board.SequenceSize
		col1 = (col1 + 1) % board.SequenceSize
		col2 = (col2 + 1) % board.SequenceSize
	}
	return b
}

func assertSampleBoard(t *testing.T, b board.Board) {
	assert.Assert(t, b.IsValid())
	assert.Assert(t, b.IsValidCell(4))
	assert.Assert(t, b.IsValidCell(40))
	assert.Assert(t, board.Equivalent(b, newSampleBoard()), "Boards are not equal", b)
}

func TestValidBoard(t *testing.T) {
	b := newSampleBoard()
	// cspell:disable-next-line
	board.Write(b, board.NewWriter(func(s string) { t.Log(strings.ReplaceAll(s, "\n", "")) }), "vrcst")
	assertSampleBoard(t, b)

	b.Set(40, board.Value(5))
	assert.Assert(t, !b.IsValid())

	assert.Assert(t, !b.IsValidCell(40))
	assert.Assert(t, b.IsValidCell(4)) // single read-only

	b.Set(40, 0)
	assertSampleBoard(t, b)
}

func TestClone(t *testing.T) {
	b := newSampleBoard()

	play := b.Clone(board.Play)
	assertSampleBoard(t, play)
	assert.Equal(t, board.Play, play.Mode())

	edit := b.Clone(board.Edit)
	assertSampleBoard(t, edit)
	assert.Equal(t, board.Edit, edit.Mode())
}

func TestPlay(t *testing.T) {
	b := newSampleBoard()
	play := b.Clone(board.Play)
	assertSampleBoard(t, play)
	assert.Equal(t, board.Play, play.Mode())

	assert.Assert(t, cmp.Panics(func() { play.SetReadOnly(55, board.Value(6)) }))
	assert.Equal(t, board.Value(0), play.Get(55))
	assertSampleBoard(t, play)

	play.Set(55, board.Value(6))
	assert.Equal(t, board.Value(6), play.Get(55))
	assert.Assert(t, play.IsValid())
	assert.Assert(t, board.ContainsReadOnly(play, b))
	assert.Assert(t, board.ContainsAll(play, b))

	play.Set(40, board.Value(5))
	assert.Equal(t, board.Value(5), play.Get(40))
	assert.Assert(t, !play.IsValid())
	assert.Assert(t, board.ContainsReadOnly(play, b))
	assert.Assert(t, board.ContainsAll(play, b))

	play.Set(40, board.Value(9))
	assert.Assert(t, play.IsValid())
	assert.Assert(t, board.ContainsReadOnly(play, b))
	assert.Assert(t, board.ContainsAll(play, b))

	play.Restart()
	assert.Equal(t, board.Value(0), play.Get(55))
	assert.Equal(t, board.Value(0), play.Get(40))
	// only the first row is read-only
	assert.Equal(t, board.FullValueSet(), play.RowSet(0))
	for row := 1; row < 9; row++ {
		assert.Equal(t, board.EmptyValueSet(), play.RowSet(row))
	}
	assert.Assert(t, board.ContainsReadOnly(b, play))
}
