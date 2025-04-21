package board

import (
	"bufio"
	"strings"
)

type boardBase struct {
	mode          BoardMode
	readOnlyFlags bitSet81
	values        [BoardSize]Value
}

func (b *boardBase) Mode() BoardMode {
	return b.mode
}

func (b *boardBase) Get(index int) Value {
	return b.values[index]
}

func (b *boardBase) IsEmpty(index int) bool {
	return b.Get(index) == 0
}

func (b *boardBase) IsReadOnly(index int) bool {
	return b.readOnlyFlags.Get(index)
}

func (b *boardBase) String() string {
	var sb strings.Builder
	WriteValues(b, bufio.NewWriter(&sb))
	return sb.String()
}

func (b *boardBase) init(mode BoardMode) {
	b.mode = mode
}

func (b *boardBase) setInternal(index int, v Value, readOnly bool) Value {
	if b.mode == Immutable {
		panic("Cannot play in immutable mode")
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

func (b *boardBase) copyValues(other *boardBase) {
	copy(b.values[:], other.values[:])
	b.readOnlyFlags = other.readOnlyFlags
}

func (b *boardBase) calcSequence(s Sequence) (vs ValueSet, dupes ValueSet) {
	for si := range s.Size() {
		index := s.Get(si)
		v := b.Get(index)
		if v == 0 {
			continue
		}
		if vs.Contains(v) {
			dupes.Add(v)
		} else {
			vs.Add(v)
		}
	}
	return vs, dupes
}
