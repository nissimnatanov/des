package indexes

const (
	RelatedSize = 20
)

func RelatedSequence(index int) Sequence {
	return Sequence{relatedCache[index].indexes[:]}
}

func RelatedSet(index int) BitSet81 {
	return relatedCache[index].mask
}

type relatedInfo struct {
	indexes [RelatedSize]int
	mask    BitSet81
}

func initRelatedIndexes() [BoardSize]relatedInfo {
	cache := [BoardSize]relatedInfo{}
	for i := range BoardSize {
		related := 0
		row := RowFromIndex(i)
		col := ColumnFromIndex(i)
		square := SquareFromIndex(i)

		rs := RowSequence(row)
		for rowIndex := range rs.Indexes() {
			if rowIndex == i {
				continue
			}
			cache[i].indexes[related] = rowIndex
			related++
			cache[i].mask.Set(rowIndex, true)
		}

		cs := ColumnSequence(col)
		for colIndex := range cs.Indexes() {
			if colIndex == i {
				continue
			}
			cache[i].indexes[related] = colIndex
			related++
			cache[i].mask.Set(colIndex, true)
		}
		ss := SquareSequence(square)
		for squareIndex := range ss.Indexes() {
			if row == RowFromIndex(squareIndex) ||
				col == ColumnFromIndex(squareIndex) {
				continue
			}
			cache[i].indexes[related] = squareIndex
			related++
			cache[i].mask.Set(squareIndex, true)
		}

		if related != RelatedSize {
			panic("expected # of related indexes to be 20")
		}
	}
	return cache
}

var relatedCache = initRelatedIndexes()
