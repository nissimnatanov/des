package indexes

import (
	"math/bits"
	"strconv"
	"strings"
)

type BitSet81 struct {
	// low bits from 0 to 63
	low uint64
	// high bits from 64 to 80 (others unused)
	high uint32
}

var MinBitSet81 = BitSet81{}
var MaxBitSet81 = BitSet81{low: 0xFFFFFFFFFFFFFFFF, high: 0x1FFFF}

func (bs BitSet81) Get(index int) bool {
	switch {
	case index < 0:
		panic("negative index")
	case index < 64:
		return (bs.low & (1 << index)) != 0
	case index < 81:
		return (bs.high & (1 << (index - 64))) != 0
	default:
		panic("index out of range")
	}
}

func (bs *BitSet81) SetTo(index int, value bool) {
	if value {
		bs.Set(index)
	} else {
		bs.Reset(index)
	}
}

func (bs *BitSet81) Reset(index int) {
	switch {
	case index < 0:
		panic("negative index")
	case index < 64:
		bs.low &= ^(1 << index)
	case index < 81:
		bs.high &= ^(1 << (index - 64))
	default:
		panic("index out of range")
	}
}

func (bs *BitSet81) Set(index int) {
	switch {
	case index < 0:
		panic("negative index")
	case index < 64:
		bs.low |= (1 << index)
	case index < 81:
		bs.high |= (1 << (index - 64))
	default:
		panic("index out of range")
	}
}

func (bs *BitSet81) ResetMask(mask BitSet81) {
	bs.low &= ^mask.low
	bs.high &= ^mask.high
}

// First returns the first index that is set to true in the BitSet81 or -1
// if all are unset.
func (bs BitSet81) First() int {
	if bs.low != 0 {
		return bits.TrailingZeros64(bs.low)
	}
	if bs.high != 0 {
		return bits.TrailingZeros32(uint32(bs.high)) + 64
	}
	return -1
}

func (bs BitSet81) Indexes(yield func(int) bool) {
	low := bs.low
	var index int
	for low != 0 {
		zeros := bits.TrailingZeros64(low)
		index += zeros
		if !yield(index) {
			return
		}
		// drop zeros and the bit we just found
		low >>= zeros + 1
		index++
	}
	high := bs.high
	index = 64
	for high != 0 {
		zeros := bits.TrailingZeros32(high)
		index += zeros
		if !yield(index) {
			return
		}
		// drop zeros and the bit we just found
		high >>= zeros + 1
		index++
	}
}

func (bs BitSet81) Complement() BitSet81 {
	return BitSet81{
		low:  ^bs.low,
		high: ^bs.high & 0x1FFFF,
	}
}

func (bs BitSet81) Intersect(other BitSet81) BitSet81 {
	return BitSet81{
		low:  bs.low & other.low,
		high: bs.high & other.high,
	}
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
