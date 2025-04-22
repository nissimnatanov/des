package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board/indexes"
	"gotest.tools/v3/assert"
)

func assertSequence(t *testing.T, s indexes.Sequence, expected [9]int) {
	for i := range s.Size() {
		assert.Equal(t, expected[i], s.Get(i), "Mismatch at %v", i)
	}
	assert.Equal(t, indexes.SequenceSize, s.Size(), "Incomplete Sequence %s", s)
}

func TestRowSequence(t *testing.T) {
	assertSequence(t, indexes.RowSequence(0), [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assertSequence(t, indexes.RowSequence(4), [9]int{36, 37, 38, 39, 40, 41, 42, 43, 44})
	assertSequence(t, indexes.RowSequence(8), [9]int{72, 73, 74, 75, 76, 77, 78, 79, 80})

	assert.Equal(t, indexes.RowSequence(7).Size(), indexes.SequenceSize)
}

func TestColSequence(t *testing.T) {
	assertSequence(t, indexes.ColumnSequence(0), [9]int{0, 9, 18, 27, 36, 45, 54, 63, 72})
	assertSequence(t, indexes.ColumnSequence(3), [9]int{3, 12, 21, 30, 39, 48, 57, 66, 75})
	assertSequence(t, indexes.ColumnSequence(8), [9]int{8, 17, 26, 35, 44, 53, 62, 71, 80})

	assert.Equal(t, indexes.ColumnSequence(2).Size(), indexes.SequenceSize)
}

func TestSquareSequence(t *testing.T) {
	assertSequence(t, indexes.SquareSequence(0), [9]int{0, 1, 2, 9, 10, 11, 18, 19, 20})
	assertSequence(t, indexes.SquareSequence(5), [9]int{33, 34, 35, 42, 43, 44, 51, 52, 53})
	assertSequence(t, indexes.SquareSequence(8), [9]int{60, 61, 62, 69, 70, 71, 78, 79, 80})

	assert.Equal(t, indexes.SquareSequence(7).Size(), indexes.SequenceSize)
}

func BenchmarkIndexFromSquare(b *testing.B) {
	for range b.N {
		for square := range indexes.SequenceSize {
			for cell := range indexes.SequenceSize {
				if indexes.IndexFromSquare(square, cell) < 0 {
					b.FailNow()
				}
			}
		}
	}
}
