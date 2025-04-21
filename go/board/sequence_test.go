package board_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board"
	"gotest.tools/v3/assert"
)

func assertSequence(t *testing.T, s board.Sequence, expected [9]int) {
	for i := range s.Size() {
		assert.Equal(t, expected[i], s.Get(i), "Mismatch at %v", i)
	}
	assert.Equal(t, board.SequenceSize, s.Size(), "Incomplete Sequence %s", s)
}

func TestRowSequence(t *testing.T) {
	assertSequence(t, board.RowSequence(0), [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assertSequence(t, board.RowSequence(4), [9]int{36, 37, 38, 39, 40, 41, 42, 43, 44})
	assertSequence(t, board.RowSequence(8), [9]int{72, 73, 74, 75, 76, 77, 78, 79, 80})

	assert.Equal(t, board.RowSequence(7).Size(), board.SequenceSize)
}

func TestColSequence(t *testing.T) {
	assertSequence(t, board.ColumnSequence(0), [9]int{0, 9, 18, 27, 36, 45, 54, 63, 72})
	assertSequence(t, board.ColumnSequence(3), [9]int{3, 12, 21, 30, 39, 48, 57, 66, 75})
	assertSequence(t, board.ColumnSequence(8), [9]int{8, 17, 26, 35, 44, 53, 62, 71, 80})

	assert.Equal(t, board.ColumnSequence(2).Size(), board.SequenceSize)
}

func TestSquareSequence(t *testing.T) {
	assertSequence(t, board.SquareSequence(0), [9]int{0, 1, 2, 9, 10, 11, 18, 19, 20})
	assertSequence(t, board.SquareSequence(5), [9]int{33, 34, 35, 42, 43, 44, 51, 52, 53})
	assertSequence(t, board.SquareSequence(8), [9]int{60, 61, 62, 69, 70, 71, 78, 79, 80})

	assert.Equal(t, board.SquareSequence(7).Size(), board.SequenceSize)
}
