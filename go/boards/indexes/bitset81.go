package indexes

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

type bitSet81Unit = byte

const (
	bitSet81BitsPerUnit = 8
	bitSet81UnitCount   = 1 + (BoardSize-1)/bitSet81BitsPerUnit
	bitSet81CacheSize   = 1 << bitSet81BitsPerUnit

	bitSet81FullUnitMask  = bitSet81Unit((1 << bitSet81BitsPerUnit) - 1)
	bitSet81FullUnitCount = BoardSize / bitSet81BitsPerUnit
	// bitSet81PartialUnitMask can be either zero (if the last unit has full mask)
	// in which case the bitSet81FullMaskUnitCount is same as bitSet81UnitCount
	bitSet81PartialUnitMask = (1 << (BoardSize % bitSet81BitsPerUnit)) - 1
)

type BitSet81 [bitSet81UnitCount]bitSet81Unit

var MinBitSet81 = BitSet81{}
var MaxBitSet81 BitSet81 = bitSet81InitMax()

func bitSet81InitMax() BitSet81 {
	max := BitSet81{}

	for ui := range bitSet81FullUnitCount {
		max[ui] = bitSet81FullUnitMask
	}
	if bitSet81PartialUnitMask != 0 {
		max[bitSet81UnitCount-1] = bitSet81PartialUnitMask
	}
	return max
}

func (bs BitSet81) Get(index int) bool {
	CheckBoardIndex(index)
	bi := bs[index/bitSet81BitsPerUnit]
	return (bi & (1 << (index % bitSet81BitsPerUnit))) != 0
}

func (bs *BitSet81) Set(index int, value bool) {
	CheckBoardIndex(index)
	if value {
		bs[index/bitSet81BitsPerUnit] |= 1 << (index % bitSet81BitsPerUnit)
	} else {
		bs[index/bitSet81BitsPerUnit] &^= 1 << (index % bitSet81BitsPerUnit)
	}
}

func (bs *BitSet81) ResetMask(mask BitSet81) {
	for ui := range bitSet81UnitCount {
		bs[ui] &= ^mask[ui]
	}
}

// First returns the first index that is set to true in the BitSet81 or -1
// if all are unset.
func (bs BitSet81) First() int {
	for ui := range bitSet81UnitCount {
		u := bs[ui]
		if u != 0 {
			return ui*bitSet81BitsPerUnit + bitSet81IndexCache[u][0]
		}
	}
	return -1
}

func (bs BitSet81) Indexes(yield func(int) bool) {
	for ui := range bitSet81UnitCount {
		u := bs[ui]
		if u == 0 {
			continue
		}
		start := ui * bitSet81BitsPerUnit
		indexes := bitSet81IndexCache[u]
		for _, index := range indexes {
			if !yield(start + index) {
				return
			}
		}
	}
}

func (bs BitSet81) Complement() BitSet81 {
	cbs := BitSet81{}
	for i := range bitSet81FullUnitCount {
		cbs[i] = (^bs[i] & bitSet81FullUnitMask)
	}
	if bitSet81PartialUnitMask != 0 {
		cbs[bitSet81UnitCount-1] = (^bs[bitSet81UnitCount-1] & bitSet81PartialUnitMask)
	}
	return cbs
}

func (bs BitSet81) Intersect(other BitSet81) BitSet81 {
	cbs := BitSet81{}
	for i := range bitSet81UnitCount {
		cbs[i] = bs[i] & other[i]
	}
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

var bitSet81IndexCache = initBitSetIndexCache()

func initBitSetIndexCache() [bitSet81CacheSize][]int {
	debug.FreeOSMemory()
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	start := ms.Alloc
	var cache [bitSet81CacheSize][]int
	for i := range cache {
		for j := range bitSet81BitsPerUnit {
			if (i & (1 << j)) != 0 {
				cache[i] = append(cache[i], j)
			}
		}
	}
	runtime.ReadMemStats(ms)
	_ = start
	fmt.Println("bitSet81IndexCache", ms.Alloc-start, "bytes")
	return cache
}
