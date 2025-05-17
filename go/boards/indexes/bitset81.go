package indexes

import (
	"strconv"
	"strings"
)

type BitSet81 [11]uint8

var MaxBitSet81 BitSet81 = [11]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}

var MinBitSet81 = BitSet81{}

func (bs BitSet81) Get(index int) bool {
	CheckBoardIndex(index)
	b := bs[index/8]
	return (b & (1 << (index % 8))) != 0
}

func (bs *BitSet81) Set(index int, value bool) {
	CheckBoardIndex(index)
	if value {
		bs[index/8] |= 1 << (index % 8)
	} else {
		bs[index/8] &^= 1 << (index % 8)
	}
}

func (bs *BitSet81) ResetMask(mask BitSet81) {
	for i := range 11 {
		bs[i] &= ^mask[i]
	}
}

func (bs BitSet81) First() (int, bool) {
	for bi := range 11 {
		b := bs[bi]
		if b != 0 {
			return bi*8 + bitSetIndexCache[b][0], true
		}
	}
	return -1, false
}

func (bs BitSet81) Indexes(yield func(int) bool) {
	for bi := range 11 {
		start := bi * 8
		indexes := bitSetIndexCache[bs[bi]]
		for _, index := range indexes {
			if !yield(start + index) {
				return
			}
		}
	}
}

func (bs BitSet81) Complement() BitSet81 {
	cbs := BitSet81{}
	for i := range 10 {
		cbs[i] = ^bs[i]
	}
	cbs[10] = ^bs[10] & 0x01
	return cbs
}

// String returns a string representation of the BitSet81,
// it is not optimized (use for troubleshooting only)!
func (bs BitSet81) String() string {
	sb := strings.Builder{}
	sb.Write([]byte{'['})
	first := true
	for index := range bs.Indexes {
		if first {
			first = false
		} else {
			sb.Write([]byte{','})
		}
		sb.WriteString(strconv.Itoa(index))
	}
	sb.Write([]byte{']'})
	return sb.String()
}

var bitSetIndexCache = initBitSetIndexCache()

func initBitSetIndexCache() [256][]int {
	var cache [256][]int
	for i := range 256 {
		cache[i] = make([]int, 0, 8)
		for j := range 8 {
			if (i & (1 << j)) != 0 {
				cache[i] = append(cache[i], j)
			}
		}
	}
	return cache
}
