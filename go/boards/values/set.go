package values

import (
	"strconv"
)

// Set represents a set of values using a bit mask, it is considered immutable.
// Only the first 9 bits are used
type Set uint16

var EmptySet = Set(0)

var FullSet = Set(0x1FF)

func Intersect(vs1 Set, vs2 Set, more ...Set) Set {
	mask := vs1 & vs2
	for _, vs := range more {
		mask &= vs
	}
	return Set(mask)
}

func Union(vs1 Set, vs2 Set, more ...Set) Set {
	mask := vs1 | vs2
	for _, vs := range more {
		mask |= vs
	}
	return Set(mask)
}

func NewSet(vs ...Value) Set {
	var mask Set
	for _, v := range vs {
		mask |= v.AsSet()
	}
	return Set(mask)
}

// Values of this set, do not modify the return slice.
func (vs Set) Values() Values {
	return setInfoCache[vs].values
}

// First value is useful when set has exactly one value and
// we want to use it as a Value.
func (vs Set) First() Value {
	if vs == 0 {
		return 0
	}
	return setInfoCache[vs].values[0]
}

func (vs Set) IsEmpty() bool {
	return vs == 0
}

func (vs Set) At(i int) Value {
	return setInfoCache[vs].values[i]
}

func (vs Set) With(other Set) Set {
	return Union(vs, other)
}

func (vs Set) Without(other Set) Set {
	return Intersect(vs, other.Complement())
}

func (vs Set) Complement() Set {
	return FullSet &^ vs
}

func (vs Set) ContainsAll(other Set) bool {
	return (vs & other) == other
}

func (vs Set) ContainsAny(other Set) bool {
	return (vs & other) != 0
}

func (vs Set) Contains(v Value) bool {
	v.Validate()
	return vs.ContainsAny(v.AsSet())
}

func (vs Set) Size() int {
	return len(setInfoCache[vs].values)
}

func (vs Set) Combined() int {
	return setInfoCache[vs].combined
}

func (vs Set) String() string {
	combined := vs.Combined()
	if combined == 0 {
		return ""
	}
	return strconv.Itoa(combined)
}

type setInfo struct {
	values   []Value
	combined int
}

func newSetInfo(mask int) setInfo {
	combined := 0
	values := []Value{}

	for v := 1; v <= 9; v++ {
		vMask := 1 << (v - 1)
		if mask&vMask != 0 {
			values = append(values, Value(v))
			combined = combined*10 + v
		}
	}
	return setInfo{values, combined}
}

func initSetInfo() []setInfo {
	var valueSets []setInfo
	for mask := range 0x1FF + 1 {
		valueSets = append(valueSets, newSetInfo(mask))
	}
	return valueSets
}

var setInfoCache = initSetInfo()
