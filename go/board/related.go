package board

const (
	RelatedSize = 20
)

func RelatedSequence(index int) *Sequence {
	return &Sequence{relatedCache[index][:]}
}

func initRelatedIndexes() [BoardSize][RelatedSize]int {
	cache := [BoardSize][RelatedSize]int{}
	for i := 0; i < BoardSize; i++ {
		related := 0
		row := RowFromIndex(i)
		col := ColumnFromIndex(i)
		square := SquareFromIndex(i)

		for ri := RowSequence(row).Iterator(); ri.Next(); {
			rowIndex := ri.Value()
			if rowIndex == i {
				continue
			}
			cache[i][related] = rowIndex
			related++
		}
		for ci := ColumnSequence(col).Iterator(); ci.Next(); {
			colIndex := ci.Value()
			if colIndex == i {
				continue
			}
			cache[i][related] = colIndex
			related++
		}
		for si := SquareSequence(square).Iterator(); si.Next(); {
			squareIndex := si.Value()
			if row == RowFromIndex(squareIndex) ||
				col == ColumnFromIndex(squareIndex) {
				continue
			}
			cache[i][related] = squareIndex
			related++
		}
		if related != RelatedSize {
			panic("expected # of related indexes to be 20")
		}
	}
	return cache
}

var relatedCache = initRelatedIndexes()
