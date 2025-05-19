package values_test

import (
	"fmt"
	"testing"

	"github.com/nissimnatanov/des/go/boards/values"
)

func TestValueProperties(t *testing.T) {
	cases := []struct {
		value    values.Value
		stringer string
	}{
		{values.Value(0), "0"},
		{values.Value(0), "0"},
		{1, "1"},
		{5, "5"},
		{9, "9"},
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
		vs         values.Set
		combined   int32
		size       int
		complement values.Set
	}{
		{values.FullSet, 123456789, 9, values.EmptySet},
		{values.EmptySet, 0, 0, values.FullSet},
		{values.Value(5).AsSet(), 5, 1, values.NewSet(1, 2, 3, 4, 6, 7, 8, 9)},
		{values.NewSet(1, 9, 5), 159, 3, values.NewSet(2, 3, 4, 6, 7, 8)},
	}
	for _, c := range cases {
		combined := c.vs.Combined()
		stringer := c.vs.String()
		size := c.vs.Size()
		isEmpty := c.vs.IsEmpty()
		complement := c.vs.Complement()
		original := complement.Complement()
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
		if original != c.vs {
			t.Errorf("Complement(%v) == %v, want %v", complement, c.vs, c.complement)
		}
	}
}

func TestValueSetValueIterator(t *testing.T) {
	cases := []struct {
		vs   values.Set
		want string
	}{
		{values.FullSet, "123456789"},
		{values.EmptySet, ""},
		{values.Value(5).AsSet(), "5"},
		{values.NewSet(1, 9, 5), "159"},
	}
	for _, c := range cases {
		var got string
		for _, v := range c.vs.Values() {
			got += v.String()
		}
		if got != c.want {
			t.Errorf("ValueIterator(%q) == %q, want %q", c.vs, got, c.want)
		}
	}
}

func TestValueSetContainsValueSet(t *testing.T) {
	cases := []struct {
		vs          values.Set
		other       values.Set
		containsAll bool
		containsAny bool
	}{
		{values.FullSet, values.Value(5).AsSet(), true, true},
		{values.FullSet, values.EmptySet, true, false},
		{values.EmptySet, values.Value(5).AsSet(), false, false},
		{values.Value(6).AsSet(), values.Value(5).AsSet(), false, false},
		{values.NewSet(1, 5), values.Value(3).AsSet(), false, false},
		{values.NewSet(9, 5, 7), values.NewSet(9, 4), false, true},
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
		vs       values.Set
		other    values.Value
		contains bool
	}{
		{values.FullSet, 5, true},
		{values.FullSet, 0, false},
		{values.EmptySet, 5, false},
		{values.Value(6).AsSet(), 5, false},
		{values.NewSet(1, 5), 3, false},
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
		vs1       values.Set
		vs2       values.Set
		vs3       values.Set
		union     values.Set
		intersect values.Set
	}{
		{values.FullSet, values.NewSet(5, 6), values.NewSet(1, 5), values.FullSet, values.Value(5).AsSet()},
		{values.EmptySet, values.NewSet(5, 6), values.NewSet(1, 5), values.NewSet(1, 5, 6), values.EmptySet},
	}

	for _, c := range cases {
		union := values.Union(c.vs1, c.vs2, c.vs3)
		intersect := values.Intersect(c.vs1, c.vs2, c.vs3)

		if union != c.union {
			t.Errorf("Union(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, union, c.union)
		}
		if intersect != c.intersect {
			t.Errorf("Intersect(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, intersect, c.intersect)
		}
	}
}
