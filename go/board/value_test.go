package board_test

import (
	"fmt"
	"testing"

	"github.com/nissimnatanov/des/go/board"
)

func TestValueProperties(t *testing.T) {
	cases := []struct {
		value    board.Value
		stringer string
	}{
		{board.Value(0), "0"},
		{board.Value(0), "0"},
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
		vs         board.ValueSet
		combined   int32
		size       int
		complement board.ValueSet
	}{
		{board.FullValueSet(), 123456789, 9, board.EmptyValueSet()},
		{board.EmptyValueSet(), 0, 0, board.FullValueSet()},
		{board.Value(5).AsSet(), 5, 1, board.NewValueSet(1, 2, 3, 4, 6, 7, 8, 9)},
		{board.NewValueSet(1, 9, 5), 159, 3, board.NewValueSet(2, 3, 4, 6, 7, 8)},
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
		vs   board.ValueSet
		want string
	}{
		{board.FullValueSet(), "123456789"},
		{board.EmptyValueSet(), ""},
		{board.Value(5).AsSet(), "5"},
		{board.NewValueSet(1, 9, 5), "159"},
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
		vs          board.ValueSet
		other       board.ValueSet
		containsAll bool
		containsAny bool
	}{
		{board.FullValueSet(), board.Value(5).AsSet(), true, true},
		{board.FullValueSet(), board.EmptyValueSet(), true, false},
		{board.EmptyValueSet(), board.Value(5).AsSet(), false, false},
		{board.Value(6).AsSet(), board.Value(5).AsSet(), false, false},
		{board.NewValueSet(1, 5), board.Value(3).AsSet(), false, false},
		{board.NewValueSet(9, 5, 7), board.NewValueSet(9, 4), false, true},
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
		vs       board.ValueSet
		other    board.Value
		contains bool
	}{
		{board.FullValueSet(), 5, true},
		{board.FullValueSet(), 0, false},
		{board.EmptyValueSet(), 5, false},
		{board.Value(6).AsSet(), 5, false},
		{board.NewValueSet(1, 5), 3, false},
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
		vs1       board.ValueSet
		vs2       board.ValueSet
		vs3       board.ValueSet
		union     board.ValueSet
		intersect board.ValueSet
	}{
		{board.FullValueSet(), board.NewValueSet(5, 6), board.NewValueSet(1, 5), board.FullValueSet(), board.Value(5).AsSet()},
		{board.EmptyValueSet(), board.NewValueSet(5, 6), board.NewValueSet(1, 5), board.NewValueSet(1, 5, 6), board.EmptyValueSet()},
	}

	for _, c := range cases {
		union := board.ValueSetUnion(c.vs1, c.vs2, c.vs3)
		intersect := board.ValueSetIntersect(c.vs1, c.vs2, c.vs3)

		if union != c.union {
			t.Errorf("Union(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, union, c.union)
		}
		if intersect != c.intersect {
			t.Errorf("Intersect(%v, %v, %v) == %v, want %v", c.vs1, c.vs2, c.vs3, intersect, c.intersect)
		}
	}
}
