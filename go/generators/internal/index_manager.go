package internal

import "github.com/nissimnatanov/des/go/boards"

type indexManager struct {
	// indexes are reserved/removed from the back!
	indexes  [boards.Size]int
	remained int
	reserved int
}

func newIndexManager() *indexManager {
	m := &indexManager{}
	for i := 0; i < boards.Size; i++ {
		m.indexes[i] = i
	}
	m.remained = boards.Size
	m.reserved = 0
	return m
}

func (im *indexManager) shuffleRemained(r *Random) {
	RandShuffle(r, im.indexes[:im.remained])
}

func (im *indexManager) shuffleRemoved(r *Random) {
	start := im.remained + im.reserved
	count := len(im.indexes) - start
	RandShuffle(r, im.indexes[start:start+count])
}

func (im *indexManager) RemainedSize() int {
	return im.remained
}

func (im *indexManager) swap(i, j int) {
	if i != j {
		im.indexes[i], im.indexes[j] = im.indexes[j], im.indexes[i]
	}
}

func (im *indexManager) SwapRemained(ri1, ri2 int) {
	if ri1 < 0 || ri1 >= im.remained || ri2 < 0 || ri2 >= im.remained {
		panic("SwapRemained received bad index")
	}
	if ri1 != ri2 {
		im.swap(ri1, ri2)
	}
}

func (im *indexManager) Reserve(n int) {
	if n <= 0 || n > im.remained {
		panic("n is out of range")
	}
	if im.reserved != 0 {
		panic("cannot reserve twice")
	}

	im.reserved = n
	im.remained -= n
}

func (im *indexManager) RemoveReserved() {
	im.reserved = 0
}

func (im *indexManager) RevertReserved() {
	im.remained += im.reserved
	im.reserved = 0
}

func (im *indexManager) Remained() []int {
	return im.indexes[:im.remained]
}

func (im *indexManager) Reserved() []int {
	return im.indexes[im.remained : im.remained+im.reserved]
}

func (im *indexManager) Removed() []int {
	return im.indexes[im.remained+im.reserved:]
}

func (im *indexManager) tryPrioritizeIndex(index int) bool {
	for ii := range im.remained {
		if im.indexes[ii] == index {
			im.SwapRemained(ii, im.remained-1)
			return true
		}
	}
	return false
}

func (im *indexManager) PrioritizeIndex(index int) {
	if !im.tryPrioritizeIndex(index) {
		panic("Index to prioritize was not found")
	}
}

func (im *indexManager) RestoreRemoved(index int) {
	for ii := im.remained + im.reserved; ii < len(im.indexes); ii++ {
		if im.indexes[ii] != index {
			continue
		}

		im.swap(im.remained, ii)
		if im.reserved > 0 {
			im.swap(im.remained+im.reserved, ii)
		}
		im.remained++
		return
	}

	panic("Index to restore from removed back to remained was not found")
}

func (im *indexManager) RemoveIndex(index int) {
	im.PrioritizeIndex(index)
	im.Reserve(1)
	im.RemoveReserved()
}

func (im *indexManager) TryRemoveIndex(index int) bool {
	if !im.tryPrioritizeIndex(index) {
		return false
	}
	im.Reserve(1)
	im.RemoveReserved()
	return true
}

func (im *indexManager) clone() *indexManager {
	clone := &indexManager{
		remained: im.remained,
		reserved: im.reserved,
	}
	copy(clone.indexes[:], im.indexes[:])
	return clone
}
