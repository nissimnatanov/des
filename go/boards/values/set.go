package values

import (
	"strconv"
)

// Set represents a set of values using a bit mask, it is considered immutable.
// Only the first 9 bits are used
type Set uint16

var EmptySet = Set(0)

var FullSet = Set(0x1FF)

func Intersect(vs1, vs2 Set) Set {
	return vs1 & vs2
}

func Intersect3(vs1, vs2, vs3 Set) Set {
	return vs1 & vs2 & vs3
}

func Union(vs1, vs2 Set) Set {
	return vs1 | vs2
}

func Union3(vs1, vs2, vs3 Set) Set {
	return vs1 | vs2 | vs3
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
	return setValues[vs]
}

// First value is useful when set has exactly one value and
// we want to use it as a Value.
func (vs Set) First() Value {
	if vs == 0 {
		return 0
	}
	return setValues[vs][0]
}

func (vs Set) IsEmpty() bool {
	return vs == 0
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
	return setSize[vs]
}

func (vs Set) Combined() int {
	return setCombined[vs]
}

func (vs Set) String() string {
	combined := vs.Combined()
	if combined == 0 {
		return ""
	}
	return strconv.Itoa(combined)
}

func initValueSet(mask int) {
	combined := 0
	values := []Value{}

	for v := 1; v <= 9; v++ {
		vMask := 1 << (v - 1)
		if mask&vMask != 0 {
			values = append(values, Value(v))
			combined = combined*10 + v
		}
	}
	setSize[mask] = len(values)
	setValues[mask] = values
	setCombined[mask] = combined
}

const setMaskCacheSize = 0x200 // 512, enough for 9 bits
var setSize [setMaskCacheSize]int
var setValues [setMaskCacheSize][]Value
var setCombined [setMaskCacheSize]int

func init() {
	for mask := range setMaskCacheSize {
		initValueSet(mask)
	}
}
