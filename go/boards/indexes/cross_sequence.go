package indexes

// RowSequenceExcludeSquare returns the 6 row indexes that are not in the square
func RowSequenceExcludeSquare(row, sq int) Sequence {
	return rowNotInSquareCache[row][sq]
}

// ColumnSequenceExcludeSquare returns the 6 column indexes that are not in the square
func ColumnSequenceExcludeSquare(col, sq int) Sequence {
	return colNotInSquareCache[col][sq]
}

// SquareSequenceExcludeRow returns the 6 square indexes that are not in the row
func SquareSequenceExcludeRow(sq, row int) Sequence {
	return squareIndexesNotInRowCache[sq][row]
}

// SquareSequenceExcludeColumn returns the 6 square indexes that are not in the column
func SquareSequenceExcludeColumn(sq, col int) Sequence {
	return squareIndexesNotInColCache[sq][col]
}

var rowNotInSquareCache = initRowNotInSquareCache()
var colNotInSquareCache = initColNotInSquareCache()
var squareIndexesNotInRowCache = initSquareIndexesNotInRowCache()
var squareIndexesNotInColCache = initSquareIndexesNotInColCache()

func initRowNotInSquareCache() [BoardSequenceSize][BoardSequenceSize]Sequence {
	var cache [BoardSequenceSize][BoardSequenceSize]Sequence
	for row := range BoardSequenceSize {
		rowSeq := RowSequence(row)
		for sq := range BoardSequenceSize {
			for _, index := range rowSeq {
				if sq == SquareFromIndex(index) {
					continue
				}
				cache[row][sq] = append(cache[row][sq], index)
			}
		}
	}
	return cache
}

func initColNotInSquareCache() [BoardSequenceSize][BoardSequenceSize]Sequence {
	var cache [BoardSequenceSize][BoardSequenceSize]Sequence
	for col := range BoardSequenceSize {
		colSeq := ColumnSequence(col)
		for sq := range BoardSequenceSize {
			for _, index := range colSeq {
				if sq == SquareFromIndex(index) {
					continue
				}
				cache[col][sq] = append(cache[col][sq], index)
			}
		}
	}
	return cache
}

func initSquareIndexesNotInRowCache() [BoardSequenceSize][BoardSequenceSize]Sequence {
	var cache [BoardSequenceSize][BoardSequenceSize]Sequence
	for sq := range BoardSequenceSize {
		sqSeq := SquareSequence(sq)
		for row := range BoardSequenceSize {
			for _, index := range sqSeq {
				if row == RowFromIndex(index) {
					continue
				}
				cache[sq][row] = append(cache[sq][row], index)
			}
		}
	}
	return cache
}

func initSquareIndexesNotInColCache() [BoardSequenceSize][BoardSequenceSize]Sequence {
	var cache [BoardSequenceSize][BoardSequenceSize]Sequence
	for sq := range BoardSequenceSize {
		sqSeq := SquareSequence(sq)
		for col := range BoardSequenceSize {
			for _, index := range sqSeq {
				if col == ColumnFromIndex(index) {
					continue
				}
				cache[sq][col] = append(cache[sq][col], index)
			}
		}
	}
	return cache
}
