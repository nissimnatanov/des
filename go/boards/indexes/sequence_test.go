package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"gotest.tools/v3/assert"
)

func assertSequence(t *testing.T, s indexes.Sequence, expected [9]int) {
	for i := range len(s) {
		assert.Equal(t, expected[i], s[i], "Mismatch at %v", i)
	}
	assert.Equal(t, indexes.BoardSequenceSize, len(s), "Incomplete Sequence %s", s)
}

func TestRowSequence(t *testing.T) {
	assertSequence(t, indexes.RowSequence(0), [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assertSequence(t, indexes.RowSequence(4), [9]int{36, 37, 38, 39, 40, 41, 42, 43, 44})
	assertSequence(t, indexes.RowSequence(8), [9]int{72, 73, 74, 75, 76, 77, 78, 79, 80})

	assert.Equal(t, len(indexes.RowSequence(7)), indexes.BoardSequenceSize)
}

func TestColSequence(t *testing.T) {
	assertSequence(t, indexes.ColumnSequence(0), [9]int{0, 9, 18, 27, 36, 45, 54, 63, 72})
	assertSequence(t, indexes.ColumnSequence(3), [9]int{3, 12, 21, 30, 39, 48, 57, 66, 75})
	assertSequence(t, indexes.ColumnSequence(8), [9]int{8, 17, 26, 35, 44, 53, 62, 71, 80})

	assert.Equal(t, len(indexes.ColumnSequence(2)), indexes.BoardSequenceSize)
}

func TestSquareSequence(t *testing.T) {
	assertSequence(t, indexes.SquareSequence(0), [9]int{0, 1, 2, 9, 10, 11, 18, 19, 20})
	assertSequence(t, indexes.SquareSequence(5), [9]int{33, 34, 35, 42, 43, 44, 51, 52, 53})
	assertSequence(t, indexes.SquareSequence(8), [9]int{60, 61, 62, 69, 70, 71, 78, 79, 80})

	assert.Equal(t, len(indexes.SquareSequence(7)), indexes.BoardSequenceSize)
}

func BenchmarkIndexFromSquare(b *testing.B) {
	for b.Loop() {
		for square := range indexes.BoardSequenceSize {
			for cell := range indexes.BoardSequenceSize {
				if indexes.IndexFromSquare(square, cell) < 0 {
					b.FailNow()
				}
			}
		}
	}
}
