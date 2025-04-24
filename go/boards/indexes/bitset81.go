package indexes

import (
	"math"
)

type BitSet81 struct {
	low uint64 // indexes < 64
	hi  uint32 // indexes >= 64 (17 bits in use only)
}

const (
	lowFullMask = math.MaxUint64 // 64 bits in use
	hiFullMask  = 0x1FFFF        // 17 bits in use
)

func (bs *BitSet81) Get(index int) bool {
	CheckBoardIndex(index)

	if index < 64 {
		val := uint64(1) << index
		return val&bs.low > 0
	}
	val := uint32(1) << (index - 64)
	return val&bs.hi > 0
}

func (bs *BitSet81) Set(index int, value bool) {
	CheckBoardIndex(index)

	if index < 64 {
		mask := uint64(1) << index
		if value {
			bs.low |= mask
		} else {
			bs.low &= lowFullMask &^ mask
		}
		return
	}
	mask := uint32(1) << (index - 64)
	if value {
		bs.hi |= mask
	} else {
		bs.hi &= hiFullMask &^ mask
	}
}

func (bs *BitSet81) AllSet() bool {
	return bs.low == lowFullMask && bs.hi == hiFullMask
}

func (bs *BitSet81) Reset() {
	bs.low = 0
	bs.hi = 0
}

func (bs *BitSet81) ResetMask(mask BitSet81) {
	bs.low &= lowFullMask &^ mask.low
	bs.hi &= hiFullMask &^ mask.hi
}

func (bs *BitSet81) SetAll(val bool) {
	if val {
		bs.low = lowFullMask
		bs.hi = hiFullMask
		return
	}
	bs.low = 0
	bs.hi = 0
}

func (bs BitSet81) Indexes(yield func(int) bool) {
	mask64 := uint64(1)
	for i := range 64 {
		if (bs.low&mask64) != 0 && !yield(i) {
			return
		}
		mask64 <<= 1
	}
	mask32 := uint32(1)
	for i := 64; i < 81; i++ {
		if (bs.hi&mask32) != 0 && !yield(i) {
			return
		}
		mask32 <<= 1
	}
}

func (bs BitSet81) Complement() BitSet81 {
	return BitSet81{
		low: ^bs.low,
		hi:  (^bs.hi & hiFullMask),
	}
}
