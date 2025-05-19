package boards

import (
	"fmt"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

func GetIntegrityChecks() bool {
	return integrityChecks
}

func SetIntegrityChecks(enabled bool) bool {
	prev := integrityChecks
	integrityChecks = enabled
	return prev
}

var integrityChecks bool

func (b *Game) checkIntegrity() {
	if !GetIntegrityChecks() {
		return
	}

	var valueCounts [10]int

	for i, v := range b.AllValues {
		valueCounts[v]++
		if v != 0 {
			// check this value is disallowed in other places
			rs := indexes.RelatedSequence(i)
			for _, related := range rs {
				rv := b.values[related]
				if rv == 0 {
					if b.AllowedValues(related).Contains(v) {
						panic("value should not be allowed")
					}
				} else if rv == v {
					// ensure one of them is marked as wrong
					if !b.IsReadOnly(related) {
						if b.IsValidCell(related) {
							panic("Dupe value not marked as invalid")
						}
					}
					if !b.IsReadOnly(i) {
						if b.IsValidCell(i) {
							panic("Dupe value not marked as invalid")
						}
					}

					if b.IsReadOnly(related) && b.IsReadOnly(i) {
						if b.IsValidCell(i) || b.IsValidCell(related) {
							panic("Dupe read-only values not marked as invalid")
						}
					}
				}
			}
			if b.AllowedValues(i) != values.EmptySet {
				panic(fmt.Sprintf(
					"allowed values for non-empty cell %v must be empty: actual %v\n%v",
					i, b.AllowedValues(i), b.String()))
			}
		} else {
			// check that disallowed values are a union of row/column/square
			disallowedValuesExpected := values.Union(
				b.relatedValues(i),
				b.allowedValues.GetDisallowedByUser(i))
			allowedSet := b.AllowedValues(i)
			if allowedSet != disallowedValuesExpected.Complement() {
				panic(fmt.Sprintf(
					"wrong allowed values for cell %v: expected %v, actual %v\n%v",
					i, disallowedValuesExpected.Complement(), b.AllowedValues(i), b.String()))
			}
		}
	}

	if valueCounts != b.valueCounts {
		panic(fmt.Sprintf(
			"wrong value counts: expected %v, actual %v\n%v",
			valueCounts, b.valueCounts, b.String()))
	}
	for si := range SequenceSize {
		rowSeqValues := b.sequenceValues(indexes.RowSequence(si))
		if b.RowValues(si) != rowSeqValues {
			panic(fmt.Sprintf(
				"wrong row values for row %d: expected %v, actual %v\n%v",
				si, rowSeqValues, b.RowValues(si), b))
		}
		colSeqValues := b.sequenceValues(indexes.ColumnSequence(si))
		if b.ColumnValues(si) != colSeqValues {
			panic(fmt.Sprintf(
				"wrong column values for column %d: expected %v, actual %v\n%v",
				si, colSeqValues, b.ColumnValues(si), b))
		}
		sqSeqValues := b.sequenceValues(indexes.SquareSequence(si))
		if b.SquareValues(si) != sqSeqValues {
			panic(fmt.Sprintf(
				"wrong square values for square %d: expected %v, actual %v\n%v",
				si, sqSeqValues, b.SquareValues(si), b))
		}
	}
}
