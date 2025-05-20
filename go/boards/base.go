package boards

import (
	"github.com/nissimnatanov/des/go/boards/indexes"
	"github.com/nissimnatanov/des/go/boards/values"
)

type Base interface {
	Get(i int) values.Value
	IsReadOnly(i int) bool
	IsValidCell(i int) bool
}

type base struct {
	values        [Size]values.Value
	readOnlyFlags indexes.BitSet81
	mode          Mode
}

// AllValues include empty and non-empty cells, for empty value is 0
func (b *base) AllValues(yield func(i int, v values.Value) bool) {
	for i, v := range b.values {
		if !yield(i, v) {
			return
		}
	}
}

func (b *base) Get(index int) values.Value {
	return b.values[index]
}

func (b *base) IsReadOnly(index int) bool {
	return b.readOnlyFlags.Get(index)
}

func (b *base) init(mode Mode) {
	b.mode = mode
}

func (b *base) setInternal(index int, v values.Value, readOnly bool) values.Value {
	if b.mode == Immutable {
		panic("Cannot play in immutable or solution mode")
	}
	if b.mode == Play && (readOnly || b.IsReadOnly(index)) {
		panic("Edit mode is not allowed")
	}

	if v == 0 && readOnly {
		panic("Empty cell cannot be read-only")
	}

	v.Validate()
	prev := b.values[index]
	b.values[index] = v
	b.readOnlyFlags.SetTo(index, readOnly)
	return prev
}

func (b *base) copyValues(other *base) {
	copy(b.values[:], other.values[:])
	b.readOnlyFlags = other.readOnlyFlags
}

func (b *base) calcSequence(s indexes.Sequence) (vs values.Set, dupes values.Set) {
	for _, index := range s {
		v := b.values[index]
		if v == 0 {
			continue
		}
		if vs.Contains(v) {
			dupes = dupes.With(v.AsSet())
		} else {
			vs = vs.With(v.AsSet())
		}
	}
	return vs, dupes
}
