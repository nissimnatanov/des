package board

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/board/indexes"
	"github.com/nissimnatanov/des/go/board/values"
)

type boardImpl struct {
	base

	rowSets              [indexes.SequenceSize]values.Set
	colSets              [indexes.SequenceSize]values.Set
	squareSets           [indexes.SequenceSize]values.Set
	valueCounts          [indexes.SequenceSize + 1]int
	validFlags           indexes.BitSet81
	userDisallowedValues [Size]values.Set

	// Allowed values cache is special - the moment at least one value changes, it is easier to invalidate
	// all indexes instead of recalculating related ones.
	allowedValuesCache           [Size]values.Set
	allowedValuesCacheValidFlags indexes.BitSet81
}

func (b *boardImpl) RowSet(row int) values.Set {
	return b.rowSets[row]
}

func (b *boardImpl) ColumnSet(col int) values.Set {
	return b.colSets[col]
}

func (b *boardImpl) SquareSet(square int) values.Set {
	return b.squareSets[square]
}

func (b *boardImpl) AllowedSet(index int) values.Set {
	v := b.Get(index)
	if v != 0 {
		return values.EmptySet()
	}

	if b.allowedValuesCacheValidFlags.Get(index) {
		return b.allowedValuesCache[index]
	}

	disallowedValues := values.Union(b.relatedValues(index), b.userDisallowedValues[index])
	allowedValues := disallowedValues.Complement()
	b.allowedValuesCacheValidFlags.Set(index, true)
	b.allowedValuesCache[index] = allowedValues
	return allowedValues
}

func (b *boardImpl) relatedValues(index int) values.Set {
	return values.Union(
		b.RowSet(indexes.RowFromIndex(index)),
		b.ColumnSet(indexes.ColumnFromIndex(index)),
		b.SquareSet(indexes.SquareFromIndex(index)))
}

func (b *boardImpl) Count(v values.Value) int {
	return b.valueCounts[v]
}

func (b *boardImpl) FreeCellCount() int {
	return b.Count(0)
}

func (b *boardImpl) IsValidCell(index int) bool {
	return b.validFlags.Get(index)
}

func (b *boardImpl) IsValid() bool {
	return b.validFlags.AllSet()
}

func (b *boardImpl) IsSolved() bool {
	return b.IsValid() && b.FreeCellCount() == 0
}

// Empty board in Edit mode
func New() Board {
	var b boardImpl
	b.init(Edit)
	b.checkIntegrity()
	return &b
}

func (b *boardImpl) Clone(mode Mode) Board {
	if mode == Immutable && b.mode == Immutable {
		return b
	}

	newBoard := &boardImpl{}
	b.cloneInto(mode, newBoard)
	return newBoard
}

func (b *boardImpl) cloneInto(mode Mode, dst *boardImpl) {
	dst.init(mode)
	dst.copyValues(&b.base)
	dst.copyStats(b)
	dst.checkIntegrity()
}

func (b *boardImpl) CloneInto(mode Mode, dst Board) {
	dstImpl, ok := dst.(*boardImpl)
	if !ok {
		panic(fmt.Errorf("Cannot CloneInto into a board of type %T", dst))
	}
	b.cloneInto(mode, dstImpl)
}

func (b *boardImpl) Set(index int, v values.Value) {
	b.setInternal(index, v, false)
}

func (b *boardImpl) SetReadOnly(index int, v values.Value) {
	b.setInternal(index, v, true)
}

func (b *boardImpl) Disallow(index int, v values.Value) {
	b.DisallowSet(index, v.AsSet())
}

func (b *boardImpl) DisallowSet(index int, vs values.Set) {
	if b.mode == Immutable {
		panic("Cannot set disallowed values on immutable board")
	}
	if vs.IsEmpty() {
		// does not make sense
		panic("Nothing to disallow")
	}

	b.userDisallowedValues[index] = values.Union(b.userDisallowedValues[index], vs)
	// if the flags are not valid, this op is useless, but also harmless
	b.allowedValuesCache[index] = b.allowedValuesCache[index].Without(vs)
}

func (b *boardImpl) DisallowReset(index int) {
	b.userDisallowedValues[index] = values.EmptySet()
	b.allowedValuesCacheValidFlags.Set(index, false)
}

func (b *boardImpl) setInternal(index int, v values.Value, readOnly bool) values.Value {
	previousValue := b.base.setInternal(index, v, readOnly)

	// stats
	needToRecalculateAll := b.updateStats(index, previousValue, v)
	if needToRecalculateAll {
		b.recalculateAllStats()
	}
	return previousValue
}

func (b *boardImpl) Restart() {
	if b.mode == Immutable {
		panic("Cannot restart an immutable board")
	}

	// for faster reset, update the values first then force recalculation of all the stats
	for i := range Size {
		if !b.IsReadOnly(i) {
			b.base.setInternal(i, 0, false)
		}
	}

	b.recalculateAllStats()
}

// only sets non-zero values
func (b *boardImpl) init(mode Mode) {
	b.base.init(mode)
	b.valueCounts[0] = Size
	b.validFlags.SetAll(true)
}

func (b *boardImpl) String() string {
	var sb strings.Builder
	WriteValues(b, bufio.NewWriter(&sb))
	return sb.String()
}

func (b *boardImpl) copyStats(source *boardImpl) {
	if source == nil {
		panic("Cannot copy nil board")
	}

	copy(b.rowSets[:], source.rowSets[:])
	copy(b.colSets[:], source.colSets[:])
	copy(b.squareSets[:], source.squareSets[:])
	copy(b.valueCounts[:], source.valueCounts[:])
	b.validFlags = source.validFlags
	copy(b.userDisallowedValues[:], source.userDisallowedValues[:])
	copy(b.allowedValuesCache[:], source.allowedValuesCache[:])
	b.allowedValuesCacheValidFlags = source.allowedValuesCacheValidFlags
}

func (b *boardImpl) updateStats(index int, oldValue, newValue values.Value) bool {
	if oldValue == newValue {
		return false
	}

	// if board was valid before and the new value does not appear in related cells,
	// there is no need to re-validate the board.
	valid := b.IsValid()
	if newValue != 0 {
		valid = valid && !b.relatedValues(index).Contains(newValue)
	}
	if !valid {
		// recalculate all - do not care about performance for this case ...
		return true
	}

	row := indexes.RowFromIndex(index)
	col := indexes.ColumnFromIndex(index)
	square := indexes.SquareFromIndex(index)

	if oldValue != 0 {
		b.rowSets[row] = b.rowSets[row].Without(oldValue.AsSet())
		b.colSets[col] = b.colSets[col].Without(oldValue.AsSet())
		b.squareSets[square] = b.squareSets[square].Without(oldValue.AsSet())
	}
	b.valueCounts[oldValue]--

	if newValue != 0 {
		b.rowSets[row] = b.rowSets[row].With(newValue.AsSet())
		b.colSets[col] = b.colSets[col].With(newValue.AsSet())
		b.squareSets[square] = b.squareSets[square].With(newValue.AsSet())
	}
	b.valueCounts[newValue]++
	if oldValue != 0 || newValue == 0 {
		// If old value was present, we cannot just remove it from the allowed set of
		// related indexes since the same value may appear in other related cells.
		// It is faster to just invalidate the allowed values cache and let them
		// be recalculated on demand.
		b.allowedValuesCacheValidFlags.ResetMask(indexes.RelatedSet(index))
		return false
	}
	// if we added new value over empty space, than it is totally safe to exclude this
	// value from the allowed values of related cells.
	relatedIndexes := indexes.RelatedSequence(index)
	for ri := range relatedIndexes.Size() {
		relatedIndex := relatedIndexes.Get(ri)
		if relatedIndex == index || !b.IsEmpty(relatedIndex) ||
			!b.allowedValuesCacheValidFlags.Get(relatedIndex) {
			continue
		}
		b.allowedValuesCache[relatedIndex] =
			b.allowedValuesCache[relatedIndex].Without(newValue.AsSet())
	}

	return false
}

func (b *boardImpl) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags.SetAll(true)
	// force recalculation of allowed values
	b.allowedValuesCacheValidFlags.Reset()

	// value counts
	for i := range indexes.SequenceSize {
		b.valueCounts[i] = 0
	}
	for i := range Size {
		b.valueCounts[b.Get(i)]++
	}

	// init rowSets, colSets; and squareSets
	// validFlags are unset if dupe detected
	for seq := range indexes.SequenceSize {
		b.rowSets[seq] = b.validateSequence(indexes.RowSequence(seq))
		b.colSets[seq] = b.validateSequence(indexes.ColumnSequence(seq))
		b.squareSets[seq] = b.validateSequence(indexes.SquareSequence(seq))
	}

	b.checkIntegrity()
}

func (b *boardImpl) validateSequence(s indexes.Sequence) values.Set {
	vs, dupes := b.calcSequence(s)
	for vi := range dupes.Size() {
		b.markSequenceInvalid(dupes.At(vi), s)
	}
	return vs
}

func (b *boardImpl) markSequenceInvalid(v values.Value, s indexes.Sequence) {
	readOnly := []int{}
	foundReadWrite := false

	for si := range s.Size() {
		index := s.Get(si)
		if b.Get(index) != v {
			continue
		}
		if b.IsReadOnly(index) {
			readOnly = append(readOnly, index)
		} else {
			foundReadWrite = true
			b.validFlags.Set(index, false)
		}
	}

	if !foundReadWrite && len(readOnly) > 1 {
		for i := range readOnly {
			b.validFlags.Set(i, false)
		}
	}
}

func (b *boardImpl) checkIntegrity() {
	if !GetIntegrityChecks() {
		return
	}

	var valueCounts [indexes.SequenceSize + 1]int

	for i := range Size {
		v := b.Get(i)
		valueCounts[v] += 1

		if v != 0 {
			// check this value is disallowed in other places
			rs := indexes.RelatedSequence(i)
			for ri := range rs.Size() {
				related := rs.Get(ri)
				rv := b.Get(related)
				if rv == 0 {
					if b.AllowedSet(related).Contains(v) {
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
		} else {
			// check that disallowed values are a union of row/column/square
			row := indexes.RowFromIndex(i)
			col := indexes.ColumnFromIndex(i)
			square := indexes.SquareFromIndex(i)

			disallowedValuesExpected := values.Union(
				b.rowSets[row],
				b.colSets[col],
				b.squareSets[square],
				b.userDisallowedValues[i])

			if b.AllowedSet(i) != disallowedValuesExpected.Complement() {
				panic(
					fmt.Sprintf(
						"wrong allowed values for cell %v: expected %v, actual %v\n%v",
						i, disallowedValuesExpected.Complement(), b.AllowedSet(i), b.String()))
			}
		}
	}

	for v := range indexes.SequenceSize {
		if valueCounts[v] != b.valueCounts[v] {
			panic(
				fmt.Sprintf(
					"wrong value counts for %v: expected %v, actual %v\n%v",
					v, valueCounts[v], b.valueCounts[v], b.String()))
		}
	}
}
