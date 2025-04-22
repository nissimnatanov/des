package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board/indexes"
	"gotest.tools/v3/assert"
)

func TestRowFromIndex(t *testing.T) {
	assert.Equal(t, 5, indexes.RowFromIndex(45))
}

func TestColFromIndex(t *testing.T) {
	assert.Equal(t, 1, indexes.ColumnFromIndex(46))
}

func TestIndexFromCoordinates(t *testing.T) {
	assert.Equal(t, 47, indexes.IndexFromCoordinates(5, 2))
}

func TestIndexFromSquare(t *testing.T) {
	assert.Equal(t, 0, indexes.IndexFromSquare(0, 0))
	assert.Equal(t, 35, indexes.IndexFromSquare(5, 2))
	assert.Equal(t, 80, indexes.IndexFromSquare(8, 8))
}

func TestSquareFromIndex(t *testing.T) {
	assert.Equal(t, 0, indexes.SquareFromIndex(2))
	assert.Equal(t, 4, indexes.SquareFromIndex(32))
	assert.Equal(t, 8, indexes.SquareFromIndex(80))
}

func TestSquareCellFromIndexTest(t *testing.T) {
	assert.Equal(t, 2, indexes.SquareCellFromIndex(2))
	assert.Equal(t, 5, indexes.SquareCellFromIndex(41))
	assert.Equal(t, 8, indexes.SquareCellFromIndex(80))
}

// more tests in sequence_test.go
