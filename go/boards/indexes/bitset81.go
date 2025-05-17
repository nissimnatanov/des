package indexes

import (
	"strconv"
	"strings"
)

type BitSet81 struct {
	bits [11]uint8
}

var maxBitSet81 = BitSet81{
	bits: [11]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
}

func (bs *BitSet81) Get(index int) bool {
	CheckBoardIndex(index)
	b := bs.bits[index/8]
	return (b & (1 << (index % 8))) != 0
}

func (bs *BitSet81) Set(index int, value bool) {
	CheckBoardIndex(index)
	if value {
		bs.bits[index/8] |= 1 << (index % 8)
	} else {
		bs.bits[index/8] &^= 1 << (index % 8)
	}
}

func (bs *BitSet81) AllSet() bool {
	return bs.bits == maxBitSet81.bits
}

func (bs *BitSet81) Reset() {
	clear(bs.bits[:])
}

func (bs *BitSet81) ResetMask(mask BitSet81) {
	for i := range 11 {
		bs.bits[i] &= ^mask.bits[i]
	}
}

func (bs *BitSet81) SetAll(val bool) {
	if !val {
		clear(bs.bits[:])
		return
	}
	bs.bits = maxBitSet81.bits
}

func (bs BitSet81) First() (int, bool) {
	for bi := range 11 {
		b := bs.bits[bi]
		if b != 0 {
			return bi*8 + bitSetIndexCache[b][0], true
		}
	}
	return -1, false
}

func (bs BitSet81) Indexes(yield func(int) bool) {
	for bi := range 11 {
		start := bi * 8
		indexes := bitSetIndexCache[bs.bits[bi]]
		for _, index := range indexes {
			if !yield(start + index) {
				return
			}
		}
	}
}

func (bs BitSet81) Complement() BitSet81 {
	bs = BitSet81{}
	for i := range 10 {
		bs.bits[i] = ^bs.bits[i]
	}
	bs.bits[10] = ^bs.bits[10] & 0x01
	return bs
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

var bitSetIndexCache [256][]int = initBitSetIndexCache()

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
