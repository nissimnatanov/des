package board

import (
	"github.com/nissimnatanov/des/go/board/indexes"
	"github.com/nissimnatanov/des/go/board/values"
)

type base struct {
	values        [Size]values.Value
	readOnlyFlags indexes.BitSet81
	mode          Mode
}

func (b *base) Mode() Mode {
	return b.mode
}

func (b *base) Get(index int) values.Value {
	return b.values[index]
}

func (b *base) IsEmpty(index int) bool {
	return b.Get(index) == 0
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
	prev := b.Get(index)
	b.values[index] = v
	b.readOnlyFlags.Set(index, readOnly)
	return prev
}

func (b *base) copyValues(other *base) {
	copy(b.values[:], other.values[:])
	b.readOnlyFlags = other.readOnlyFlags
}

func (b *base) calcSequence(s indexes.Sequence) (vs values.Set, dupes values.Set) {
	for si := range s.Size() {
		index := s.Get(si)
		v := b.Get(index)
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
