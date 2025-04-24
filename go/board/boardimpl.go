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

	freeCellCount        int
	validFlags           indexes.BitSet81
	userDisallowedValues [Size]values.Set

	// AllowedValuesCache must be always up-2-date, it does not include used disallowed values
	// (user disallowed values are excluded from the public AllowedSet API).
	// If cell is empty, it contains the allowed values for the cell.
	// If cell is not empty, it is empty.
	allowedValuesCache [Size]values.Set
}

// TODO: we may remove RowSet, ColumnSet, SquareSet if not needed or leave for non-fast solver
func (b *boardImpl) RowSet(row int) values.Set {
	return b.sequenceValues(indexes.RowSequence(row))
}

func (b *boardImpl) ColumnSet(col int) values.Set {
	return b.sequenceValues(indexes.ColumnSequence(col))
}

func (b *boardImpl) SquareSet(square int) values.Set {
	return b.sequenceValues(indexes.SquareSequence(square))
}

func (b *boardImpl) AllowedSet(index int) values.Set {
	return b.allowedValuesCache[index].Without(b.userDisallowedValues[index])
}

func (b *boardImpl) relatedValues(index int) values.Set {
	return b.sequenceValues(indexes.RelatedSequence(index))
}

func (b *boardImpl) sequenceValues(s indexes.Sequence) values.Set {
	values := values.EmptySet()
	for i := range s.Indexes() {
		v := b.Get(i)
		if v != 0 {
			values = values.With(v.AsSet())
		}
	}
	return values
}

func (b *boardImpl) FreeCellCount() int {
	return b.freeCellCount
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
	b.initZeroStats(Edit)
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
	dst.base.init(mode)
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
}

func (b *boardImpl) DisallowReset(index int) {
	b.userDisallowedValues[index] = values.EmptySet()
}

func (b *boardImpl) setInternal(index int, v values.Value, readOnly bool) values.Value {
	previousValue := b.base.setInternal(index, v, readOnly)
	b.updateStats(index, previousValue, v)
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
func (b *boardImpl) initZeroStats(mode Mode) {
	b.base.init(mode)
	b.freeCellCount = Size
	for i := range b.allowedValuesCache {
		b.allowedValuesCache[i] = values.FullSet()
	}
	b.validFlags.SetAll(true)
	b.checkIntegrity()
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

	b.freeCellCount = source.freeCellCount
	b.validFlags = source.validFlags
	copy(b.userDisallowedValues[:], source.userDisallowedValues[:])
	copy(b.allowedValuesCache[:], source.allowedValuesCache[:])
}

func (b *boardImpl) updateStats(index int, oldValue, newValue values.Value) {
	if oldValue == newValue {
		return
	}

	// if board was valid before and the new value does not appear in related cells,
	// there is no need to re-validate the board.
	isValid := b.IsValid()
	if isValid && newValue != 0 {
		var allowedValues values.Set
		if oldValue != 0 {
			// If cell had a value before, its allowed values cache was 0, hence we cannot use it
			// to validate the new value. Instead, recalculate it based on the related cells.
			allowedValues = b.relatedValues(index).Complement()
		} else {
			// if cell was empty before, its allowed values cache was valid
			allowedValues = b.allowedValuesCache[index]
		}
		isValid = isValid && allowedValues.Contains(newValue)
	}
	if !isValid {
		// Recalculate the whole board to update Valid state on the cells. Since the built-in
		// solvers should never hit it, do not care about performance for this case.
		b.recalculateAllStats()
		return
	}

	if oldValue == 0 {
		b.freeCellCount--
	}
	if newValue == 0 {
		b.freeCellCount++
		// if we set non-empty to empty, recalculate the allowed values
		b.allowedValuesCache[index] = b.relatedValues(index).Complement()
	} else {
		b.allowedValuesCache[index] = values.EmptySet()
	}

	relatedSeq := indexes.RelatedSequence(index)
	for relatedIndex := range relatedSeq.Indexes() {
		if relatedIndex == index || !b.IsEmpty(relatedIndex) {
			continue
		}

		if oldValue != 0 {
			// If old value was present, we cannot just add it to the allowed set of
			// related indexes since the same value may appear in other related cells.
			b.allowedValuesCache[relatedIndex] = b.relatedValues(relatedIndex).Complement()
		}
		if newValue != 0 {
			// if we added new value than it is totally safe to exclude this
			// value from the allowed values of related cells.
			b.allowedValuesCache[relatedIndex] =
				b.allowedValuesCache[relatedIndex].Without(newValue.AsSet())
		}
	}
	b.checkIntegrity()
}

func (b *boardImpl) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags.SetAll(true)

	// value counts
	b.freeCellCount = 0
	for i := range Size {
		if b.IsEmpty(i) {
			b.freeCellCount++
			b.allowedValuesCache[i] = b.relatedValues(i).Complement()
		} else {
			b.allowedValuesCache[i] = values.EmptySet()
		}
	}

	// init rowSets, colSets; and squareSets
	// validFlags are unset if dupe detected
	for seq := range SequenceSize {
		b.validateSequence(indexes.RowSequence(seq))
		b.validateSequence(indexes.ColumnSequence(seq))
		b.validateSequence(indexes.SquareSequence(seq))
	}

	b.checkIntegrity()
}

func (b *boardImpl) validateSequence(s indexes.Sequence) {
	_, dupes := b.calcSequence(s)
	for v := range dupes.Values() {
		b.markSequenceInvalid(v, s)
	}
}

func (b *boardImpl) markSequenceInvalid(v values.Value, s indexes.Sequence) {
	readOnly := []int{}
	foundReadWrite := false

	for index := range s.Indexes() {
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

	var freeCellCount int

	for i := range Size {
		v := b.Get(i)
		if v == 0 {
			freeCellCount++
		}

		if v != 0 {
			// check this value is disallowed in other places
			rs := indexes.RelatedSequence(i)
			for related := range rs.Indexes() {
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
			if b.AllowedSet(i) != values.EmptySet() {
				panic(fmt.Sprintf(
					"allowed values for non-empty cell %v must be empty: actual %v\n%v",
					i, b.AllowedSet(i), b.String()))
			}
		} else {
			// check that disallowed values are a union of row/column/square
			disallowedValuesExpected := values.Union(
				b.relatedValues(i),
				b.userDisallowedValues[i])

			if b.AllowedSet(i) != disallowedValuesExpected.Complement() {
				panic(fmt.Sprintf(
					"wrong allowed values for cell %v: expected %v, actual %v\n%v",
					i, disallowedValuesExpected.Complement(), b.AllowedSet(i), b.String()))
			}
		}
	}

	if freeCellCount != b.freeCellCount {
		panic(fmt.Sprintf(
			"wrong free cell counts: expected %v, actual %v\n%v",
			freeCellCount, b.freeCellCount, b.String()))
	}
}
