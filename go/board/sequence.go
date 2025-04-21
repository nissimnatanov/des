package board

type Sequence struct {
	indexes []int
}

func (s Sequence) Get(i int) int {
	return s.indexes[i]
}

func (s Sequence) Size() int {
	return len(s.indexes)
}
