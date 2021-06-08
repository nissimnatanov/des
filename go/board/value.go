package board

import "fmt"

type Value int8

const (
	Empty Value = iota // 0
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	MinValue = One
	MaxValue = Nine
)

type ValueSet struct {
	mask int16
}

func EmptySet() ValueSet {
	return ValueSet{0}
}

func FullSet() ValueSet {
	return ValueSet{0x1FF}
}

func (v Value) IsEmpty() bool {
	return v == Empty
}

func (v Value) Validate() {
	if v < 0 || v > 9 {
		panic("Wrong Value")
	}
}

func (vs ValueSet) IsEmpty() bool {
	return vs.mask == 0
}

func (vs *ValueSet) Add(v Value) {
	v.Validate()
	*vs = Union(*vs, v.AsSet())
}

func (vs *ValueSet) Remove(v Value) {
	v.Validate()
	*vs = Intersect(*vs, v.AsSet().Complement())
}

func (v Value) AsSet() ValueSet {
	v.Validate()
	if v.IsEmpty() {
		return EmptySet()
	}
	return ValueSet{1 << (v - 1)}
}

func Union(vs1 ValueSet, vs2 ValueSet, more ...ValueSet) ValueSet {
	mask := vs1.mask | vs2.mask
	for _, vs := range more {
		mask |= vs.mask
	}
	return ValueSet{mask}
}

func NewValueSet(v1 Value, v2 Value, more ...Value) ValueSet {
	mask := v1.AsSet().mask | v2.AsSet().mask
	for _, v := range more {
		mask |= v.AsSet().mask
	}
	return ValueSet{mask}
}

func Intersect(vs1 ValueSet, vs2 ValueSet, more ...ValueSet) ValueSet {
	mask := vs1.mask & vs2.mask
	for _, vs := range more {
		mask &= vs.mask
	}
	return ValueSet{mask}
}

func (vs ValueSet) Complement() ValueSet {
	return ValueSet{FullSet().mask &^ vs.mask}
}

func (vs ValueSet) ContainsAll(other ValueSet) bool {
	return (vs.mask & other.mask) == other.mask
}

func (vs ValueSet) ContainsAny(other ValueSet) bool {
	return (vs.mask & other.mask) != 0
}

func (vs ValueSet) Contains(v Value) bool {
	v.Validate()
	return vs.ContainsAny(v.AsSet())
}

func (vs ValueSet) Size() int {
	return len(valueSetInfoCache[vs.mask].values)
}

func (vs ValueSet) Combined() int {
	return valueSetInfoCache[vs.mask].combined
}

// Use AS IS - implementation is already a reference
type ValueIterator struct {
	values []Value
	cur    int
}

func (vs ValueSet) Iterator() *ValueIterator {
	return &ValueIterator{valueSetInfoCache[vs.mask].values, -1}
}

func (vi *ValueIterator) Next() bool {
	vi.cur++
	return vi.cur < len(vi.values)
}

func (vi *ValueIterator) Value() Value {
	return vi.values[vi.cur]
}

func (v Value) String() string {
	return fmt.Sprint(int8(v))
}

func (vs ValueSet) String() string {
	combined := vs.Combined()
	if combined == 0 {
		return ""
	}
	return fmt.Sprint(combined)
}

type valueSetInfo struct {
	values   []Value
	combined int
}

func newSetInfo(mask int) valueSetInfo {
	combined := 0
	values := []Value{}

	for v := 1; v <= 9; v++ {
		vmask := 1 << (v - 1)
		if mask&vmask != 0 {
			values = append(values, Value(v))
			combined = combined*10 + v
		}
	}
	return valueSetInfo{values, combined}
}

func initValueSetInfo() []valueSetInfo {
	var valueSets []valueSetInfo
	for mask := 0; mask <= 0x1FF; mask++ {
		valueSets = append(valueSets, newSetInfo(mask))
	}
	return valueSets
}

var valueSetInfoCache = initValueSetInfo()
