package indexes

import (
	"iter"
	"slices"
)

type Sequence struct {
	indexes []int
}

func (s Sequence) Indexes() iter.Seq[int] {
	return slices.Values(s.indexes)
}

func (s Sequence) At(i int) int {
	return s.indexes[i]
}

func (s Sequence) Size() int {
	return len(s.indexes)
}
