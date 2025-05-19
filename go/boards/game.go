package boards

import (
	"bufio"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type Game struct {
	base

	valueCounts [10]int // zero for empty cells
	validFlags  indexes.BitSet81

	// allowedValues must be always up-2-date
	allowedValues values.Allowed

	rowSets    [SequenceSize]values.Set
	colSets    [SequenceSize]values.Set
	squareSets [SequenceSize]values.Set
}

func (b *Game) Mode() Mode {
	return b.mode
}

func (b *Game) IsEmpty(index int) bool {
	return b.values[index] == 0
}

func (b *Game) RowValues(row int) values.Set {
	return b.rowSets[row]
}

func (b *Game) ColumnValues(col int) values.Set {
	return b.colSets[col]
}

func (b *Game) SquareValues(sq int) values.Set {
	return b.squareSets[sq]
}

// Hint01 returns the first cell that has either zero or one allowed value.
// If no such cell exists, it returns -1 and false.
func (b *Game) Hint01() int {
	return b.allowedValues.Hint01()
}

func (b *Game) AllAllowedValues(yield func(int, values.Set) bool) {
	for i := range Size {
		if b.values[i] == 0 && !yield(i, b.allowedValues.Get(i)) {
			return
		}
	}
}

func (b *Game) AllowedValues(index int) values.Set {
	return b.allowedValues.Get(index)
}

func (b *Game) AllowedValuesIn(seq indexes.Sequence) func(yield func(int, values.Set) bool) {
	return func(yield func(int, values.Set) bool) {
		for _, i := range seq {
			if b.values[i] == 0 && !yield(i, b.allowedValues.Get(i)) {
				return
			}
		}
	}
}

func (b *Game) relatedValues(index int) values.Set {
	return b.sequenceValues(indexes.RelatedSequence(index))
}

func (b *Game) sequenceValues(seq indexes.Sequence) values.Set {
	values := values.EmptySet
	for _, i := range seq {
		v := b.values[i]
		if v != 0 {
			values = values.With(v.AsSet())
		}
	}
	return values
}

func (b *Game) FreeCellCount() int {
	return b.valueCounts[0]
}

// if v is empty, returns the number of empty cells
func (b *Game) ValueCount(v values.Value) int {
	return b.valueCounts[v]
}

func (b *Game) IsValidCell(index int) bool {
	return b.validFlags.Get(index)
}

func (b *Game) IsValid() bool {
	return b.validFlags == indexes.MaxBitSet81
}

func (b *Game) IsSolved() bool {
	return b.IsValid() && b.FreeCellCount() == 0
}

// Empty board in Edit mode
func New() *Game {
	b := &Game{}
	b.initZeroStats(Edit)
	return b
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

func (b *Game) DisallowValue(index int, v values.Value) {
	b.DisallowValues(index, v.AsSet())
}

func (b *Game) DisallowValues(index int, vs values.Set) {
	if b.mode == Immutable {
		panic("Cannot set disallowed values on immutable board")
	}
	if vs.IsEmpty() {
		// does not make sense
		panic("Nothing to disallow")
	}
	b.allowedValues.DisallowByUser(index, vs)
}

func (b *Game) ResetDisallowedByUser(index int) {
	if b.mode == Immutable {
		panic("Cannot reset disallowed values on immutable board")
	}
	b.allowedValues.ResetDisallowedByUser(index)
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
	clear(b.valueCounts[:])
	b.valueCounts[0] = Size
	b.allowedValues.AllowAll()
	b.validFlags = indexes.MaxBitSet81
	clear(b.rowSets[:])
	clear(b.colSets[:])
	clear(b.squareSets[:])
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

	b.valueCounts = source.valueCounts
	b.validFlags = source.validFlags
	b.allowedValues = source.allowedValues.Clone()
	b.rowSets = source.rowSets
	b.colSets = source.colSets
	b.squareSets = source.squareSets
}

func (b *Game) updateStats(index int, oldValue, newValue values.Value) {
	if oldValue == newValue {
		return
	}

	// if board was valid before and the new value does not appear in related cells,
	// there is no need to re-validate the board.
	isValid := b.IsValid()
	var currentRelatedValues values.Set
	var hasCurrentRelatedValues bool
	if isValid && newValue != 0 {
		var allowedValues values.Set
		if oldValue != 0 {
			// If cell had a value before, its allowed values cache was 0, hence we cannot use it
			// to validate the new value. Instead, recalculate it based on the related cells.
			currentRelatedValues = b.relatedValues(index)
			hasCurrentRelatedValues = true
			allowedValues = currentRelatedValues.Complement()
		} else {
			// if cell was empty before, its allowed values cache was valid
			allowedValues = b.allowedValues.GetByRelated(index)
		}
		isValid = isValid && allowedValues.Contains(newValue)
	}
	if !isValid {
		// Recalculate the whole board to update Valid state on the cells. Since the built-in
		// solvers should never hit it, do not care about performance for this case.
		b.recalculateAllStats()
		return
	}

	b.valueCounts[oldValue]--
	b.valueCounts[newValue]++

	row := indexes.RowFromIndex(index)
	col := indexes.ColumnFromIndex(index)
	sq := indexes.SquareFromIndex(index)

	if oldValue != 0 {
		b.rowSets[row] = b.rowSets[row].Without(oldValue.AsSet())
		b.colSets[col] = b.colSets[col].Without(oldValue.AsSet())
		b.squareSets[sq] = b.squareSets[sq].Without(oldValue.AsSet())
	}

	if newValue == 0 {
		// if we set non-empty to empty, recalculate the allowed values
		if !hasCurrentRelatedValues {
			currentRelatedValues = b.relatedValues(index)
			hasCurrentRelatedValues = true
		}
		b.allowedValues.ReportEmpty(index, currentRelatedValues)
	} else {
		b.allowedValues.ReportPresent(index)
		b.rowSets[row] = b.rowSets[row].With(newValue.AsSet())
		b.colSets[col] = b.colSets[col].With(newValue.AsSet())
		b.squareSets[sq] = b.squareSets[sq].With(newValue.AsSet())
	}
	if oldValue == 0 && newValue != 0 {
		// optimization - disallowing related indexes for a new value runs 2% faster
		// if run directly against values.Allowed.
		b.allowedValues.DisallowRelatedOf(index, newValue)
	} else {
		relatedSeq := indexes.RelatedSequence(index)
		for _, relatedIndex := range relatedSeq {
			if relatedIndex == index || b.values[relatedIndex] != 0 {
				continue
			}

			if oldValue != 0 {
				// If old value was present, we cannot just add it to the allowed set of
				// related indexes since the same value may appear in other related cells.
				b.allowedValues.ReportEmpty(relatedIndex, b.relatedValues(relatedIndex))
			}
			if newValue != 0 {
				// if we added new value than it is totally safe to include this
				// value to the disallowed values based on the related cells.
				b.allowedValues.DisallowRelated(relatedIndex, newValue)
			}
		}
	}
	b.checkIntegrity()
}

func (b *Game) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags = indexes.MaxBitSet81

	// value counts
	clear(b.valueCounts[:])
	for i := range Size {
		b.valueCounts[b.values[i]]++
		if b.values[i] == 0 {
			b.allowedValues.ReportEmpty(i, b.relatedValues(i))
		} else {
			b.allowedValues.ReportPresent(i)
		}
	}

	// init rowSets, colSets; and squareSets
	// validFlags are unset if dupe detected
	for si := range SequenceSize {
		b.rowSets[si] = b.processSequence(indexes.RowSequence(si))
		b.colSets[si] = b.processSequence(indexes.ColumnSequence(si))
		b.squareSets[si] = b.processSequence(indexes.SquareSequence(si))
	}

	b.checkIntegrity()
}

func (b *Game) processSequence(s indexes.Sequence) values.Set {
	vs, dupes := b.calcSequence(s)
	for _, v := range dupes.Values() {
		b.markSequenceInvalid(v, s)
	}
	return vs
}

func (b *Game) markSequenceInvalid(v values.Value, s indexes.Sequence) {
	readOnly := []int{}
	foundReadWrite := false

	for _, index := range s {
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
		for _, roi := range readOnly {
			b.validFlags.Set(roi, false)
		}
	}
}
