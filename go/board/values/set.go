package values

import (
	"strconv"
)

func EmptySet() Set {
	return Set{0}
}

func FullSet() Set {
	return Set{0x1FF}
}

func Intersect(vs1 Set, vs2 Set, more ...Set) Set {
	mask := vs1.mask & vs2.mask
	for _, vs := range more {
		mask &= vs.mask
	}
	return Set{mask}
}

func Union(vs1 Set, vs2 Set, more ...Set) Set {
	mask := vs1.mask | vs2.mask
	for _, vs := range more {
		mask |= vs.mask
	}
	return Set{mask}
}

func NewSet(vs ...Value) Set {
	var mask int16
	for _, v := range vs {
		mask |= v.AsSet().mask
	}
	return Set{mask}
}

// Set represents a set of values using a bit mask, it is considered immutable.
type Set struct {
	mask int16
}

func (vs Set) IsEmpty() bool {
	return vs.mask == 0
}

func (vs Set) At(i int) Value {
	return setInfoCache[vs.mask].values[i]
}

func (vs Set) With(other Set) Set {
	return Union(vs, other)
}

func (vs Set) Without(other Set) Set {
	return Intersect(vs, other.Complement())
}

func (vs Set) Complement() Set {
	return Set{FullSet().mask &^ vs.mask}
}

func (vs Set) ContainsAll(other Set) bool {
	return (vs.mask & other.mask) == other.mask
}

func (vs Set) ContainsAny(other Set) bool {
	return (vs.mask & other.mask) != 0
}

func (vs Set) Contains(v Value) bool {
	v.Validate()
	return vs.ContainsAny(v.AsSet())
}

func (vs Set) Size() int {
	return len(setInfoCache[vs.mask].values)
}

func (vs Set) Combined() int {
	return setInfoCache[vs.mask].combined
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
