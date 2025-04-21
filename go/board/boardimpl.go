package board

import "fmt"

type boardImpl struct {
	boardBase

	rowSets              [SequenceSize]ValueSet
	colSets              [SequenceSize]ValueSet
	squareSets           [SequenceSize]ValueSet
	valueCounts          [SequenceSize + 1]int
	validFlags           bitSet81
	userDisallowedValues [BoardSize]ValueSet

	// Allowed values cache is special - the moment at least one value changes, it is easier to invalidate
	// all indexes instead of recalculating related ones.
	allowedValuesCache           [BoardSize]ValueSet
	allowedValuesCacheValidFlags bitSet81
}

func (b *boardImpl) RowSet(row int) ValueSet {
	return b.rowSets[row]
}

func (b *boardImpl) ColumnSet(col int) ValueSet {
	return b.colSets[col]
}

func (b *boardImpl) SquareSet(square int) ValueSet {
	return b.squareSets[square]
}

func (b *boardImpl) AllowedSet(index int) ValueSet {
	v := b.Get(index)
	if v != 0 {
		return EmptyValueSet()
	}

	if b.allowedValuesCacheValidFlags.Get(index) {
		return b.allowedValuesCache[index]
	}

	disallowedValues := ValueSetUnion(b.relatedSet(index), b.userDisallowedValues[index])
	allowedValues := disallowedValues.Complement()
	b.allowedValuesCacheValidFlags.Set(index, true)
	b.allowedValuesCache[index] = allowedValues
	return allowedValues
}

func (b *boardImpl) relatedSet(index int) ValueSet {
	return ValueSetUnion(
		b.RowSet(RowFromIndex(index)),
		b.ColumnSet(ColumnFromIndex(index)),
		b.SquareSet(SquareFromIndex(index)))
}

func (b *boardImpl) Count(v Value) int {
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
	return &b
}

func (b *boardImpl) Clone(mode BoardMode) Board {
	if mode == Immutable && b.mode == Immutable {
		return b
	}

	var newBoard boardImpl
	newBoard.init(mode)
	newBoard.copyValues(&b.boardBase)
	newBoard.copyStats(b)
	return &newBoard
}

func (b *boardImpl) Set(index int, v Value) {
	b.setInternal(index, v, false)
}

func (b *boardImpl) SetReadOnly(index int, v Value) {
	b.setInternal(index, v, true)
}

func (b *boardImpl) Disallow(index int, v Value) {
	b.DisallowSet(index, v.AsSet())
}

func (b *boardImpl) DisallowSet(index int, vs ValueSet) {
	if b.mode == Immutable {
		panic("Cannot set disallowed values on immutable board")
	}
	if vs.IsEmpty() {
		// does not make sense
		panic("Nothing to disallow")
	}

	b.userDisallowedValues[index] = ValueSetUnion(b.userDisallowedValues[index], vs)
	b.allowedValuesCacheValidFlags.Set(index, false)
}

func (b *boardImpl) DisallowReset(index int) {
	b.userDisallowedValues[index] = EmptyValueSet()
	b.allowedValuesCacheValidFlags.Set(index, false)
}

func (b *boardImpl) setInternal(index int, v Value, readOnly bool) Value {
	previousValue := b.boardBase.setInternal(index, v, readOnly)

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
	for i := range BoardSize {
		if !b.IsReadOnly(i) {
			b.setInternal(i, 0, false)
		}
	}

	b.recalculateAllStats()
}

// only sets non-zero values
func (b *boardImpl) init(mode BoardMode) {
	b.boardBase.init(mode)
	b.valueCounts[0] = BoardSize
	b.validFlags.SetAll(true)
	b.checkIntegrity()
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

func (b *boardImpl) updateStats(index int, oldValue, newValue Value) bool {
	if oldValue == newValue {
		return false
	}

	// if board was valid before and the new value does not appear in related cells,
	// there is no need to re-validate the board.
	valid := b.IsValid()
	if newValue != 0 {
		valid = valid && !b.relatedSet(index).Contains(newValue)
	}
	if !valid {
		// recalculate all - do not care about performance for this case ...
		return true
	}

	row := RowFromIndex(index)
	col := ColumnFromIndex(index)
	square := SquareFromIndex(index)

	if oldValue != 0 {
		b.rowSets[row].Remove(oldValue)
		b.colSets[col].Remove(oldValue)
		b.squareSets[square].Remove(oldValue)
	}
	b.valueCounts[oldValue]--

	if newValue != 0 {
		b.rowSets[row].Add(newValue)
		b.colSets[col].Add(newValue)
		b.squareSets[square].Add(newValue)
	}
	b.valueCounts[newValue]++
	b.allowedValuesCacheValidFlags.Reset()
	return false
}

func (b *boardImpl) recalculateAllStats() {
	// assume valid unless proven otherwise inside calcSequenceSet
	b.validFlags.SetAll(true)
	// force recalculation of allowed values
	b.allowedValuesCacheValidFlags.Reset()

	// value counts
	for i := range SequenceSize {
		b.valueCounts[i] = 0
	}
	for i := range BoardSize {
		b.valueCounts[b.Get(i)]++
	}

	// init rowSets, colSets; and squareSets
	// validFlags are unset if dupe detected
	for seq := range SequenceSize {
		b.rowSets[seq] = b.validateSequence(RowSequence(seq))
		b.colSets[seq] = b.validateSequence(ColumnSequence(seq))
		b.squareSets[seq] = b.validateSequence(SquareSequence(seq))
	}
}

func (b *boardImpl) validateSequence(s Sequence) ValueSet {
	vs, dupes := b.calcSequence(s)
	for vi := dupes.Iterator(); vi.Next(); {
		b.markSequenceInvalid(vi.Value(), s)
	}
	return vs
}

func (b *boardImpl) markSequenceInvalid(v Value, s Sequence) {
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

	var valueCounts [SequenceSize + 1]int

	for i := range BoardSize {
		v := b.Get(i)
		valueCounts[v] += 1

		if v != 0 {
			// check this value is disallowed in other places
			rs := RelatedSequence(i)
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
			row := RowFromIndex(i)
			col := ColumnFromIndex(i)
			square := SquareFromIndex(i)

			disallowedValuesExpected := ValueSetUnion(
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

	for v := range SequenceSize {
		if valueCounts[v] != b.valueCounts[v] {
			panic(
				fmt.Sprintf(
					"wrong value counts for %v: expected %v, actual %v\n%v",
					v, valueCounts[v], b.valueCounts[v], b.String()))
		}
	}
}
