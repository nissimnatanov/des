package boards

import (
	"encoding/json"
	"strings"

	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type Game struct {
	base

	valueCounts [10]int // zero for empty cells
	validFlags  indexes.BitSet81

	rowSets    [SequenceSize]values.Set
	colSets    [SequenceSize]values.Set
	squareSets [SequenceSize]values.Set

	disallowedByRelated [indexes.BoardSize]values.Set
	disallowedByUser    [indexes.BoardSize]values.Set
	allowed             [indexes.BoardSize]values.Set
	emptyCells          indexes.BitSet81
	hints01             indexes.BitSet81
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
	return b.hints01.First()
}

func (b *Game) Hints01() indexes.BitSet81 {
	return b.hints01
}

func (b *Game) AllAllowedValues(yield func(int, values.Set) bool) {
	for i := range Size {
		if b.values[i] == 0 && !yield(i, b.AllowedValues(i)) {
			return
		}
	}
}

func (b *Game) AllowedValues(index int) values.Set {
	return b.allowed[index]
}

func (b *Game) AllowedValuesIn(seq indexes.Sequence) func(yield func(int, values.Set) bool) {
	return func(yield func(int, values.Set) bool) {
		for _, i := range seq {
			if b.values[i] == 0 && !yield(i, b.AllowedValues(i)) {
				return
			}
		}
	}
}

func (b *Game) EmptyCells() indexes.BitSet81 {
	return b.emptyCells
}

// calcRelatedValues includes the value of the cell itself, if not empty
// and the values of all related cells (row, column, square).
func (b *Game) calcRelatedValues(index int) values.Set {
	return values.Union3(
		b.rowSets[indexes.RowFromIndex(index)],
		b.colSets[indexes.ColumnFromIndex(index)],
		b.squareSets[indexes.SquareFromIndex(index)])
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

func (b *Game) isValidCell(index int) bool {
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

func (b *Game) Reset(indexes ...int) {
	for _, index := range indexes {
		b.setInternal(index, 0, false)
	}
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
	b.disallowedByUser[index] = b.disallowedByUser[index].With(vs)
	b.updateAllowed(index)
}

func (b *Game) DisallowedByUser(index int) values.Set {
	return b.disallowedByUser[index]
}

func (b *Game) getDisallowedByUser(index int) values.Set {
	return b.disallowedByUser[index]
}

func (b *Game) ResetDisallowedByUser(index int) {
	if b.mode == Immutable {
		panic("Cannot reset disallowed values on immutable board")
	}
	b.disallowedByUser[index] = values.EmptySet
	b.updateAllowed(index)
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
	b.validFlags = indexes.MaxBitSet81
	clear(b.rowSets[:])
	clear(b.colSets[:])
	clear(b.squareSets[:])
	clear(b.disallowedByRelated[:])
	clear(b.disallowedByUser[:])
	for i := range Size {
		b.allowed[i] = values.FullSet
	}
	b.emptyCells = indexes.MaxBitSet81
	b.hints01 = indexes.MinBitSet81
	b.checkIntegrity()
}

func (b *Game) String() string {
	var sb strings.Builder
	writeValues(b, &sb)
	return sb.String()
}

func (b *Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(Serialize(b))
}

func (b *Game) copyStats(source *Game) {
	if source == nil {
		panic("Cannot copy nil board")
	}

	b.valueCounts = source.valueCounts
	b.validFlags = source.validFlags
	b.rowSets = source.rowSets
	b.colSets = source.colSets
	b.squareSets = source.squareSets
	b.disallowedByRelated = source.disallowedByRelated
	b.disallowedByUser = source.disallowedByUser
	b.allowed = source.allowed
	b.emptyCells = source.emptyCells
	b.hints01 = source.hints01
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
			// If cell had a value before, its allowed values cache was 0, hence we cannot use
			// it to validate the new value. Instead, recalculate it based on the related cells.
			allowedValues = b.calcRelatedValues(index).Complement()
		} else {
			// if cell was empty before, its allowed values cache was valid
			allowedValues = b.valuesAllowedByRelated(index)
		}
		isValid = isValid && allowedValues.Contains(newValue)
	}
	if !isValid {
		// Recalculate the whole board to update Valid state on the cells. Since the built-in
		// solvers should never hit it, do not care about performance for this case.
		b.recalculateAllStats()
		return
	}

	// If board was valid and remains valid with a new value, we can optimize the
	// recalculation of the allowed values by touching only the affected cells.
	// All calculations below are based on the fact that the board is valid.

	row := indexes.RowFromIndex(index)
	col := indexes.ColumnFromIndex(index)
	sq := indexes.SquareFromIndex(index)

	// first, update the row/col/square sets - related values calculation depends on it
	if oldValue != 0 {
		// since board was valid, we can safely assume that the oldValue was present
		// only once per row/col/square
		oldValueSet := oldValue.AsSet()
		b.rowSets[row] = b.rowSets[row].Without(oldValueSet)
		b.colSets[col] = b.colSets[col].Without(oldValueSet)
		b.squareSets[sq] = b.squareSets[sq].Without(oldValueSet)
	}
	if newValue != 0 {
		newValueSet := newValue.AsSet()
		b.rowSets[row] = b.rowSets[row].With(newValueSet)
		b.colSets[col] = b.colSets[col].With(newValueSet)
		b.squareSets[sq] = b.squareSets[sq].With(newValueSet)
	}

	b.valueCounts[oldValue]--
	b.valueCounts[newValue]++

	if newValue == 0 {
		// if we set non-empty to empty, recalculate the allowed values
		b.updateOnEmpty(index, b.calcRelatedValues(index))
	} else {
		b.updateOnPresent(index)
	}

	if oldValue == 0 && newValue != 0 {
		// optimization - disallowing related indexes for a new value runs 2% faster
		// if run directly against values.Allowed.
		b.disallowRelatedOf(index, newValue)
	} else {
		relatedSeq := indexes.RelatedSequence(index)
		for _, relatedIndex := range relatedSeq {
			if relatedIndex == index || b.values[relatedIndex] != 0 {
				continue
			}

			if oldValue != 0 {
				// If old value was present, we cannot just add it to the allowed set of
				// related indexes since the same value may appear in other related cells.
				b.updateOnEmpty(relatedIndex, b.calcRelatedValues(relatedIndex))
			}
			if newValue != 0 {
				// if we added new value than it is totally safe to include this
				// value to the disallowed values based on the related cells.
				b.disallowRelated(relatedIndex, newValue)
			}
		}
	}
	b.checkIntegrity()
}

func (b *Game) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags = indexes.MaxBitSet81

	// init rowSets, colSets; and squareSets
	// validFlags are unset if dupe detected
	for si := range SequenceSize {
		b.rowSets[si] = b.processSequence(indexes.RowSequence(si))
		b.colSets[si] = b.processSequence(indexes.ColumnSequence(si))
		b.squareSets[si] = b.processSequence(indexes.SquareSequence(si))
	}

	clear(b.valueCounts[:])
	for i := range Size {
		b.valueCounts[b.values[i]]++
		if b.values[i] == 0 {
			b.updateOnEmpty(i, b.calcRelatedValues(i))
		} else {
			b.updateOnPresent(i)
		}
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
			b.validFlags.Reset(index)
		}
	}

	if !foundReadWrite && len(readOnly) > 1 {
		for _, roi := range readOnly {
			b.validFlags.Reset(roi)
		}
	}
}

func (b *Game) valuesAllowedByRelated(index int) values.Set {
	return b.disallowedByRelated[index].Complement()
}

func (b *Game) disallowRelated(index int, v values.Value) {
	if !b.emptyCells.Get(index) {
		panic("disallowing a value in a cell that has a value")
	}
	vs := v.AsSet()
	b.disallowedByRelated[index] = b.disallowedByRelated[index].With(vs)
	b.updateAllowed(index)
}

// ReportPresent is used when board cell has a value set
func (b *Game) updateOnPresent(index int) {
	b.disallowedByRelated[index] = values.FullSet
	b.disallowedByUser[index] = values.EmptySet
	b.allowed[index] = values.EmptySet
	b.emptyCells.Reset(index)
	b.hints01.Reset(index)
}

func (b *Game) updateOnEmpty(index int, related values.Set) {
	b.disallowedByRelated[index] = related
	b.emptyCells.Set(index)
	b.updateAllowed(index)
}

func (b *Game) disallowRelatedOf(index int, newValue values.Value) {
	newValueSet := newValue.AsSet()
	relatedEmpty := b.emptyCells.Intersect(indexes.RelatedSet(index))
	for related := range relatedEmpty.Indexes {
		b.disallowedByRelated[related] = b.disallowedByRelated[related].With(newValueSet)
		b.updateAllowed(related)
	}
}

func (b *Game) updateAllowed(index int) {
	allowed := values.Union(b.disallowedByRelated[index], b.disallowedByUser[index]).Complement()
	b.allowed[index] = allowed
	isHint01 := allowed.Size() <= 1
	b.hints01.SetTo(index, isHint01)
}
