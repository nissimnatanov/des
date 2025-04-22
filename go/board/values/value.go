package values

import "fmt"

type Value int8

func (v Value) Validate() {
	if v < 0 || v > 9 {
		panic("Value out of range")
	}
}

func (v Value) AsSet() Set {
	v.Validate()
	if v == 0 {
		return EmptySet()
	}
	return Set{1 << (v - 1)}
}

func (v Value) String() string {
	return fmt.Sprint(int8(v))
}
