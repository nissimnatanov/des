package board

const (
	RelatedSize = 20
)

func RelatedSequence(index int) Sequence {
	return Sequence{relatedCache[index][:]}
}

func initRelatedIndexes() [BoardSize][RelatedSize]int {
	cache := [BoardSize][RelatedSize]int{}
	for i := range BoardSize {
		related := 0
		row := RowFromIndex(i)
		col := ColumnFromIndex(i)
		square := SquareFromIndex(i)
		rs := RowSequence(row)

		for ri := range rs.Size() {
			rowIndex := rs.Get(ri)
			if rowIndex == i {
				continue
			}
			cache[i][related] = rowIndex
			related++
		}

		cs := ColumnSequence(col)
		for ci := range cs.Size() {
			colIndex := cs.Get(ci)
			if colIndex == i {
				continue
			}
			cache[i][related] = colIndex
			related++
		}
		ss := SquareSequence(square)
		for si := range ss.Size() {
			squareIndex := ss.Get(si)
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
