package indexes_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"gotest.tools/v3/assert"
)

func TestRowSequenceExcludeSquare(t *testing.T) {
	assert.DeepEqual(t,
		indexes.RowSequenceExcludeSquare(0, 0),
		indexes.Sequence{3, 4, 5, 6, 7, 8})
	assert.DeepEqual(t,
		indexes.RowSequenceExcludeSquare(8, 8),
		indexes.Sequence{72, 73, 74, 75, 76, 77})
	assert.DeepEqual(t,
		indexes.RowSequenceExcludeSquare(4, 4),
		indexes.Sequence{36, 37, 38, 42, 43, 44})
	assert.DeepEqual(t,
		indexes.RowSequenceExcludeSquare(7, 7),
		indexes.Sequence{63, 64, 65, 69, 70, 71})

	// row sequences returned as is since row 6 does not share indexes with square 1
	assert.DeepEqual(t,
		indexes.RowSequenceExcludeSquare(6, 1),
		indexes.Sequence{54, 55, 56, 57, 58, 59, 60, 61, 62})

	for row := range 9 {
		for sq := range 9 {
			rowSeqExcludeSquare := indexes.RowSequenceExcludeSquare(row, sq)
			for _, index := range rowSeqExcludeSquare {
				assert.Assert(t, indexes.RowFromIndex(index) == row)
				assert.Assert(t, indexes.SquareFromIndex(index) != sq)
			}
			assert.Assert(t, len(rowSeqExcludeSquare) == 6 || len(rowSeqExcludeSquare) == 9)
		}
	}
}

func TestColumnSequenceExcludeSquare(t *testing.T) {
	assert.DeepEqual(t,
		indexes.ColumnSequenceExcludeSquare(0, 0),
		indexes.Sequence{27, 36, 45, 54, 63, 72})
	assert.DeepEqual(t,
		indexes.ColumnSequenceExcludeSquare(8, 8),
		indexes.Sequence{8, 17, 26, 35, 44, 53})

	// column sequences returned as is since column 1 does not share indexes with square 8
	assert.DeepEqual(t,
		indexes.ColumnSequenceExcludeSquare(1, 8),
		indexes.Sequence{1, 10, 19, 28, 37, 46, 55, 64, 73})

	for col := range 9 {
		for sq := range 9 {
			colSeqExcludeSquare := indexes.ColumnSequenceExcludeSquare(col, sq)
			for _, index := range colSeqExcludeSquare {
				assert.Assert(t, indexes.ColumnFromIndex(index) == col)
				assert.Assert(t, indexes.SquareFromIndex(index) != sq)
			}
			assert.Assert(t, len(colSeqExcludeSquare) == 6 || len(colSeqExcludeSquare) == 9)
		}
	}
}

func TestSquareSequenceExcludeRow(t *testing.T) {
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeRow(0, 0),
		indexes.Sequence{9, 10, 11, 18, 19, 20})
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeRow(8, 8),
		indexes.Sequence{60, 61, 62, 69, 70, 71})

	// square sequences returned as is since square 1 does not share indexes with row 8
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeRow(1, 8),
		indexes.Sequence{3, 4, 5, 12, 13, 14, 21, 22, 23})

	for sq := range 9 {
		for row := range 9 {
			sqSeqExcludeRow := indexes.SquareSequenceExcludeRow(sq, row)
			for _, index := range sqSeqExcludeRow {
				assert.Assert(t, indexes.SquareFromIndex(index) == sq)
				assert.Assert(t, indexes.RowFromIndex(index) != row)
			}
			assert.Assert(t, len(sqSeqExcludeRow) == 6 || len(sqSeqExcludeRow) == 9)
		}
	}
}

func TestSquareSequenceExcludeColumn(t *testing.T) {
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeColumn(0, 0),
		indexes.Sequence{1, 2, 10, 11, 19, 20})
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeColumn(8, 8),
		indexes.Sequence{60, 61, 69, 70, 78, 79})

	// square sequences returned as is since square 1 does not share indexes with column 8
	assert.DeepEqual(t,
		indexes.SquareSequenceExcludeColumn(1, 8),
		indexes.Sequence{3, 4, 5, 12, 13, 14, 21, 22, 23})

	for sq := range 9 {
		for col := range 9 {
			sqSeqExcludeCol := indexes.SquareSequenceExcludeColumn(sq, col)
			for _, index := range sqSeqExcludeCol {
				assert.Assert(t, indexes.SquareFromIndex(index) == sq)
				assert.Assert(t, indexes.ColumnFromIndex(index) != col)
			}
			assert.Assert(t, len(sqSeqExcludeCol) == 6 || len(sqSeqExcludeCol) == 9)
		}
	}
}
