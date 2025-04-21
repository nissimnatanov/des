package board_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/board"
	"gotest.tools/v3/assert"
)

func TestRowFromIndex(t *testing.T) {
	assert.Equal(t, 5, board.RowFromIndex(45))
}

func TestColFromIndex(t *testing.T) {
	assert.Equal(t, 1, board.ColumnFromIndex(46))
}

func TestIndexFromCoordinates(t *testing.T) {
	assert.Equal(t, 47, board.IndexFromCoordinates(5, 2))
}

func TestIndexFromSquare(t *testing.T) {
	assert.Equal(t, 0, board.IndexFromSquare(0, 0))
	assert.Equal(t, 35, board.IndexFromSquare(5, 2))
	assert.Equal(t, 80, board.IndexFromSquare(8, 8))
}

func TestSquareFromIndex(t *testing.T) {
	assert.Equal(t, 0, board.SquareFromIndex(2))
	assert.Equal(t, 4, board.SquareFromIndex(32))
	assert.Equal(t, 8, board.SquareFromIndex(80))
}

func TestSquareCellFromIndexTest(t *testing.T) {
	assert.Equal(t, 2, board.SquareCellFromIndex(2))
	assert.Equal(t, 5, board.SquareCellFromIndex(41))
	assert.Equal(t, 8, board.SquareCellFromIndex(80))
}

// more tests in sequence_test.go
