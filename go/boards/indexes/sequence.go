package indexes

type Sequence struct {
	indexes []int
}

func (s Sequence) Indexes(yield func(int) bool) {
	for _, i := range s.indexes {
		if !yield(i) {
			return
		}
	}
}

func (s Sequence) At(i int) int {
	return s.indexes[i]
}

func (s Sequence) Size() int {
	return len(s.indexes)
}
