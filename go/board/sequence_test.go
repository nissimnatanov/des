package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertSequence(t *testing.T, s *Sequence, expected [9]int) {
	i := 0
	s.ForEach(func(index int) {
		assert.Equal(t, expected[i], index, "Mismatch at %v", i)
		i++
	})
	assert.Equal(t, SequenceSize, i, "Incomplete iterator at ", i)
}

func TestRowSequence(t *testing.T) {
	assertSequence(t, RowSequence(0), [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8})
	assertSequence(t, RowSequence(4), [9]int{36, 37, 38, 39, 40, 41, 42, 43, 44})
	assertSequence(t, RowSequence(8), [9]int{72, 73, 74, 75, 76, 77, 78, 79, 80})

	assert.Equal(t, RowSequence(7).Size(), SequenceSize)
}

func TestColSequence(t *testing.T) {
	assertSequence(t, ColumnSequence(0), [9]int{0, 9, 18, 27, 36, 45, 54, 63, 72})
	assertSequence(t, ColumnSequence(3), [9]int{3, 12, 21, 30, 39, 48, 57, 66, 75})
	assertSequence(t, ColumnSequence(8), [9]int{8, 17, 26, 35, 44, 53, 62, 71, 80})

	assert.Equal(t, ColumnSequence(2).Size(), SequenceSize)
}

func TestSquareSequence(t *testing.T) {
	assertSequence(t, SquareSequence(0), [9]int{0, 1, 2, 9, 10, 11, 18, 19, 20})
	assertSequence(t, SquareSequence(5), [9]int{33, 34, 35, 42, 43, 44, 51, 52, 53})
	assertSequence(t, SquareSequence(8), [9]int{60, 61, 62, 69, 70, 71, 78, 79, 80})

	assert.Equal(t, SquareSequence(7).Size(), SequenceSize)
}
