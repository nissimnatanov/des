package board

import "math"

type bitSet81 struct {
	low uint64 // indexes < 64
	hi  uint32 // indexes >= 64 (17 bits in use only)
}

const (
	lowFullMask = math.MaxUint64 // 64 bits in use
	hiFullMask  = 0x1FFFF        // 17 bits in use
)

func (bs *bitSet81) Get(index int) bool {
	CheckBoardIndex(index)

	if index < 64 {
		val := uint64(1) << index
		return val&bs.low > 0
	}
	val := uint32(1) << (index - 64)
	return val&bs.hi > 0
}

func (bs *bitSet81) Set(index int, value bool) {
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

func (bs *bitSet81) AllSet() bool {
	return bs.low == lowFullMask && bs.hi == hiFullMask
}

func (bs *bitSet81) Reset() {
	bs.low = 0
	bs.hi = 0
}

func (bs *bitSet81) SetAll(val bool) {
	if val {
		bs.low = lowFullMask
		bs.hi = hiFullMask
		return
	}
	bs.low = 0
	bs.hi = 0
}
