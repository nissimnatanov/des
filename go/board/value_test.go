package board

import (
	"fmt"
	"testing"
)

func TestValueProperties(t *testing.T) {
	cases := []struct {
		value    Value
		stringer string
	}{
		{Empty, "0"},
		{Value(0), "0"},
		{One, "1"},
		{Five, "5"},
		{Nine, "9"},
	}
	for _, c := range cases {
		got := c.value.String()
		if got != c.stringer {
			t.Errorf("String(%q) == %q, want %q", c.value, got, c.stringer)
		}
	}
}

func TestValueSetProperties(t *testing.T) {
	cases := []struct {
		vs         ValueSet
		combined   int32
		size       int
		complement ValueSet
	}{
		{FullSet(), 123456789, 9, EmptySet()},
		{EmptySet(), 0, 0, FullSet()},
		{Five.AsSet(), 5, 1, NewValueSet(One, Two, Three, Four, Six, Seven, Eight, Nine)},
		{NewValueSet(One, Nine, Five), 159, 3, NewValueSet(Two, Three, Four, Six, Seven, Eight)},
	}
	for _, c := range cases {
		combined := c.vs.Combined()
		stringer := c.vs.String()
		size := c.vs.Size()
		isEmpty := c.vs.IsEmpty()
		complement := c.vs.Complement()
		originial := complement.Complement()
		if c.combined != 0 {
			if stringer != fmt.Sprint(c.combined) {
				t.Errorf("String(%#v) == %v, want %v", c.vs, stringer, c.combined)
			}
		} else {
			if len(stringer) > 0 {
				t.Errorf("String(%#v) == %v, want empty", c.vs, stringer)
			}
		}
		if int32(combined) != c.combined {
			t.Errorf("Combined(%v) == %v, want %v", c.vs, combined, c.combined)
		}
		if size != c.size {
			t.Errorf("Size(%v) == %v, want %v", c.vs, size, c.size)
		}
		if isEmpty != (c.size == 0) {
			t.Errorf("IsEmpty(%v) == %v, want %v", c.vs, isEmpty, c.size == 0)
		}
		if complement != c.complement {
			t.Errorf("Complement(%v) == %v, want %v", c.vs, complement, c.complement)
		}
		if originial != c.vs {
			t.Errorf("Complement(%v) == %v, want %v", complement, c.vs, c.complement)
		}
	}
}

func TestValueSetValueIterator(t *testing.T) {
	cases := []struct {
		vs   ValueSet
		want string
	}{
		{FullSet(), "123456789"},
		{EmptySet(), ""},
		{Five.AsSet(), "5"},
		{NewValueSet(One, Nine, Five), "159"},
	}
	for _, c := range cases {
		var got string
		values := c.vs.Iterator()
		for values.Next() {
			got += values.Value().String()
		}
		if got != c.want {
			t.Errorf("ValueIterator(%q) == %q, want %q", c.vs, got, c.want)
		}
	}
}

func TestValueSetContainsValueSet(t *testing.T) {
	cases := []struct {
		vs          ValueSet
		other       ValueSet
		containsAll bool
		containsAny bool
	}{
		{FullSet(), Five.AsSet(), true, true},
		{FullSet(), EmptySet(), true, false},
		{EmptySet(), Five.AsSet(), false, false},
		{Six.AsSet(), Five.AsSet(), false, false},
		{NewValueSet(One, Five), Three.AsSet(), false, false},
		{NewValueSet(Nine, Five, Seven), NewValueSet(Nine, Four), false, true},
	}

	for _, c := range cases {
		containsAll := c.vs.ContainsAll(c.other)
		containsAny := c.vs.ContainsAny(c.other)

		if containsAll != c.containsAll {
			t.Errorf("(%q).ContainsAll(%q) == %t, want %t", c.vs, c.other, containsAll, c.containsAll)
		}
		if containsAny != c.containsAny {
			t.Errorf("(%q).ContainsAny(%q) == %v, want %v", c.vs, c.other, containsAny, c.containsAny)
		}
	}
}

func TestValueSetContainsValue(t *testing.T) {
	cases := []struct {
		vs       ValueSet
		other    Value
		contains bool
	}{
		{FullSet(), Five, true},
		{FullSet(), Empty, false},
		{EmptySet(), Five, false},
		{Six.AsSet(), Five, false},
		{NewValueSet(One, Five), Three, false},
	}

	for _, c := range cases {
		contains := c.vs.Contains(c.other)

		if contains != c.contains {
			t.Errorf("(%q).Contains(%q) == %t, want %t", c.vs, c.other, contains, c.contains)
		}
	}
}

func TestValueSetSetOperations(t *testing.T) {
	cases := []struct {
		vs1       ValueSet
		vs2       ValueSet
		vs3       ValueSet
		union     ValueSet
		intersect ValueSet
	}{
		{FullSet(), NewValueSet(Five, Six), NewValueSet(One, Five), FullSet(), Five.AsSet()},
		{EmptySet(), NewValueSet(Five, Six), NewValueSet(One, Five), NewValueSet(One, Five, Six), EmptySet()},
	}

	for _, c := range cases {
		union := Union(c.vs1, c.vs2, c.vs3)
		intersect := Intersect(c.vs1, c.vs2, c.vs3)

		if union != c.union {
			t.Errorf("Union(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, union, c.union)
		}
		if intersect != c.intersect {
			t.Errorf("Intersect(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, intersect, c.intersect)
		}
	}
}
