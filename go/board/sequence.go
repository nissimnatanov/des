package board

type Sequence struct {
	indexes []int
}

type SequenceIterator struct {
	indexes []int
	cur     int
}

func (s *Sequence) Get(i int) int {
	return s.indexes[i]
}

func (s *Sequence) Size() int {
	return len(s.indexes)
}

func (s *Sequence) ForEach(op func(index int)) {
	for _, index := range s.indexes {
		op(index)
	}
}

func (s *Sequence) Iterator() *SequenceIterator {
	return &SequenceIterator{s.indexes, -1}
}

func (si *SequenceIterator) Next() bool {
	si.cur++
	return si.cur < len(si.indexes)
}

// Must call Next beforehand and get true
func (si *SequenceIterator) Value() int {
	return si.indexes[si.cur]
}
