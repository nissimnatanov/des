package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowFromIndex(t *testing.T) {
	assert.Equal(t, 5, RowFromIndex(45))
}

func TestColFromIndex(t *testing.T) {
	assert.Equal(t, 1, ColumnFromIndex(46))
}

func TestIndexFromCoordinates(t *testing.T) {
	assert.Equal(t, 47, IndexFromCoordinates(5, 2))
}

func TestIndexFromSquare(t *testing.T) {
	assert.Equal(t, 0, IndexFromSquare(0, 0))
	assert.Equal(t, 35, IndexFromSquare(5, 2))
	assert.Equal(t, 80, IndexFromSquare(8, 8))
}

func TestSquareFromIndex(t *testing.T) {
	assert.Equal(t, 0, SquareFromIndex(2))
	assert.Equal(t, 4, SquareFromIndex(32))
	assert.Equal(t, 8, SquareFromIndex(80))
}

func TestSquareCellFromIndexTest(t *testing.T) {
	assert.Equal(t, 2, SquareCellFromIndex(2))
	assert.Equal(t, 5, SquareCellFromIndex(41))
	assert.Equal(t, 8, SquareCellFromIndex(80))
}

// more tests in sequence_test.go
