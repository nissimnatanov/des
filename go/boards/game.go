package boards

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type disallowedValues struct {
	byRelated values.Set
	byUser    values.Set
}

func (d disallowedValues) Complement() values.Set {
	return values.Union(d.byRelated, d.byUser).Complement()
}

type Game struct {
	base

	freeCellCount int
	validFlags    indexes.BitSet81

	// disallowedValues must be always up-2-date
	// If cell is empty, disallowedValues contains the values not allowed for the cell.
	// If cell is not empty, it is full.
	disallowedValues [Size]disallowedValues
}

func (b *Game) Mode() Mode {
	return b.mode
}

func (b *Game) IsEmpty(index int) bool {
	return b.values[index] == 0
}

func (b *Game) AllowedSets(yield func(int, values.Set) bool) {
	for i, a := range b.disallowedValues[:] {
		if b.values[i] == 0 && !yield(i, a.Complement()) {
			return
		}
	}
}

// TODO: we may remove RowSet, ColumnSet, SquareSet if not needed or leave for non-fast solver
func (b *Game) RowSet(row int) values.Set {
	return b.sequenceValues(indexes.RowSequence(row))
}

func (b *Game) ColumnSet(col int) values.Set {
	return b.sequenceValues(indexes.ColumnSequence(col))
}

func (b *Game) SquareSet(square int) values.Set {
	return b.sequenceValues(indexes.SquareSequence(square))
}

func (b *Game) AllowedSet(index int) values.Set {
	return b.disallowedValues[index].Complement()
}

func (b *Game) relatedValues(index int) values.Set {
	return b.sequenceValues(indexes.RelatedSequence(index))
}

func (b *Game) sequenceValues(s indexes.Sequence) values.Set {
	values := values.EmptySet()
	for i := range s.Indexes {
		v := b.values[i]
		if v != 0 {
			values = values.With(v.AsSet())
		}
	}
	return values
}

func (b *Game) FreeCellCount() int {
	return b.freeCellCount
}

func (b *Game) IsValidCell(index int) bool {
	return b.validFlags.Get(index)
}

func (b *Game) IsValid() bool {
	return b.validFlags.AllSet()
}

func (b *Game) IsSolved() bool {
	return b.IsValid() && b.FreeCellCount() == 0
}

// Empty board in Edit mode
func New() *Game {
	var b Game
	b.initZeroStats(Edit)
	return &b
}

func (b *Game) Clone(mode Mode) *Game {
	if mode == Immutable && b.mode == Immutable {
		return b
	}

	newBoard := &Game{}
	b.cloneInto(mode, newBoard)
	return newBoard
}

func (b *Game) cloneInto(mode Mode, dst *Game) {
	dst.init(mode)
	dst.copyValues(&b.base)
	dst.copyStats(b)
	dst.checkIntegrity()
}

func (b *Game) CloneInto(mode Mode, dst *Game) {
	b.cloneInto(mode, dst)
}

func (b *Game) Set(index int, v values.Value) {
	b.setInternal(index, v, false)
}

func (b *Game) SetReadOnly(index int, v values.Value) {
	b.setInternal(index, v, true)
}

func (b *Game) Disallow(index int, v values.Value) {
	b.DisallowSet(index, v.AsSet())
}

func (b *Game) DisallowSet(index int, vs values.Set) {
	if b.mode == Immutable {
		panic("Cannot set disallowed values on immutable board")
	}
	if vs.IsEmpty() {
		// does not make sense
		panic("Nothing to disallow")
	}

	b.disallowedValues[index].byUser = values.Union(b.disallowedValues[index].byUser, vs)
}

func (b *Game) DisallowReset(index int) {
	b.disallowedValues[index].byUser = values.EmptySet()
}

func (b *Game) setInternal(index int, v values.Value, readOnly bool) values.Value {
	previousValue := b.base.setInternal(index, v, readOnly)
	b.updateStats(index, previousValue, v)
	return previousValue
}

func (b *Game) Restart() {
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
func (b *Game) initZeroStats(mode Mode) {
	b.init(mode)
	b.freeCellCount = Size
	clear(b.disallowedValues[:])
	b.validFlags.SetAll(true)
	b.checkIntegrity()
}

func (b *Game) String() string {
	var sb strings.Builder
	WriteValues(b, bufio.NewWriter(&sb))
	return sb.String()
}

func (b *Game) copyStats(source *Game) {
	if source == nil {
		panic("Cannot copy nil board")
	}

	b.freeCellCount = source.freeCellCount
	b.validFlags = source.validFlags
	copy(b.disallowedValues[:], source.disallowedValues[:])
}

func (b *Game) updateStats(index int, oldValue, newValue values.Value) {
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
			allowedValues = b.disallowedValues[index].byRelated.Complement()
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
		b.disallowedValues[index].byRelated = b.relatedValues(index)
	} else {
		b.disallowedValues[index].byRelated = values.FullSet()
		b.disallowedValues[index].byUser = values.EmptySet()
	}

	relatedSeq := indexes.RelatedSequence(index)
	for relatedIndex := range relatedSeq.Indexes {
		if relatedIndex == index || b.values[relatedIndex] != 0 {
			continue
		}

		if oldValue != 0 {
			// If old value was present, we cannot just add it to the allowed set of
			// related indexes since the same value may appear in other related cells.
			b.disallowedValues[relatedIndex].byRelated = b.relatedValues(relatedIndex)
		}
		if newValue != 0 {
			// if we added new value than it is totally safe to include this
			// value to the disallowed values based on the related cells.
			b.disallowedValues[relatedIndex].byRelated =
				b.disallowedValues[relatedIndex].byRelated.With(newValue.AsSet())
		}
	}
	b.checkIntegrity()
}

func (b *Game) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags.SetAll(true)

	// value counts
	b.freeCellCount = 0
	for i := range Size {
		if b.values[i] == 0 {
			b.freeCellCount++
			b.disallowedValues[i].byRelated = b.relatedValues(i)
		} else {
			b.disallowedValues[i].byRelated = values.FullSet()
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

func (b *Game) validateSequence(s indexes.Sequence) {
	_, dupes := b.calcSequence(s)
	for v := range dupes.Values {
		b.markSequenceInvalid(v, s)
	}
}

func (b *Game) markSequenceInvalid(v values.Value, s indexes.Sequence) {
	readOnly := []int{}
	foundReadWrite := false

	for index := range s.Indexes {
		if b.values[index] != v {
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

func (b *Game) checkIntegrity() {
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
			for related := range rs.Indexes {
				rv := b.values[related]
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
				b.disallowedValues[i].byUser)

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
