package indexes

const BoardSequenceSize = 9
const BoardSize = 81

func CheckBoardIndex(index int) {
	if index < 0 || index >= BoardSize {
		panic("Index out of range")
	}
}

func RowFromIndex(index int) int {
	return rowFromIndexCache[index]
}

func ColumnFromIndex(index int) int {
	return colFromIndexCache[index]
}

func SquareFromIndex(index int) int {
	return squareFromIndexCache[index]
}

func SquareCellFromIndex(index int) int {
	return squareCellFromIndexCache[index]
}

func IndexFromCoordinates(row, col int) int {
	return indexFromCoordinatesCache[row][col]
}

func IndexFromSquare(square, cell int) int {
	// math here is ~4x expensive than memory access
	return indexFromSquareCache[square][cell]
}

func RowSequence(row int) Sequence {
	return Sequence{indexFromCoordinatesCache[row][:]}
}

func ColumnSequence(col int) Sequence {
	return Sequence{columnIndexes[col][:]}
}

func SquareSequence(square int) Sequence {
	return Sequence{indexFromSquareCache[square][:]}
}

func initRowFromIndex() [BoardSize]int {
	var cache [BoardSize]int
	for i := range BoardSize {
		cache[i] = i / BoardSequenceSize
	}
	return cache
}

func initColFromIndex() [BoardSize]int {
	var cache [BoardSize]int
	for i := range BoardSize {
		cache[i] = i % BoardSequenceSize
	}
	return cache
}

func initSquareFromIndex() [BoardSize]int {
	var cache [BoardSize]int
	for i := range BoardSize {
		square := i / 3
		square = (square/9)*3 + square%3
		cache[i] = square
	}
	return cache
}

func initSquareCellFromIndex() [BoardSize]int {
	var cache [BoardSize]int
	for i := range BoardSize {
		// rows (3,4,5) and (6,7,8) are equivalent to (0,1,2)
		row := (i / 9) % 3
		squareCell := i%3 + row*3
		cache[i] = squareCell
	}
	return cache
}

func initIndexFromCoordinates() [BoardSequenceSize][BoardSequenceSize]int {
	var cache [BoardSequenceSize][BoardSequenceSize]int
	for row := range BoardSequenceSize {
		for col := range BoardSequenceSize {
			cache[row][col] = row*9 + col
		}
	}
	return cache
}

func initColumnIndexes() [BoardSequenceSize][BoardSequenceSize]int {
	var cache [BoardSequenceSize][BoardSequenceSize]int
	for col := range BoardSequenceSize {
		for row := range BoardSequenceSize {
			cache[col][row] = row*9 + col
		}
	}
	return cache
}

func initIndexFromSquare() [BoardSequenceSize][BoardSequenceSize]int {
	var cache [BoardSequenceSize][BoardSequenceSize]int
	for square := range BoardSequenceSize {
		for cell := range BoardSequenceSize {
			// index of the first cell
			index := (square/3)*27 + (square%3)*3
			// add cell location relative to first one
			index += (cell/3)*9 + (cell % 3)
			cache[square][cell] = index
		}
	}
	return cache
}

var rowFromIndexCache = initRowFromIndex()
var colFromIndexCache = initColFromIndex()
var squareFromIndexCache = initSquareFromIndex()
var squareCellFromIndexCache = initSquareCellFromIndex()

// row -> list of col indexes in this row
var indexFromCoordinatesCache = initIndexFromCoordinates()

// col -> list of row indexes in this column
var columnIndexes = initColumnIndexes()

// square -> list of cell indexes in this square
var indexFromSquareCache = initIndexFromSquare()
