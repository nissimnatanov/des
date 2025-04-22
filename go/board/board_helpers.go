package board

func Equivalent(b1, b2 Board) bool {
	for i := range Size {
		if b1.Get(i) != b2.Get(i) {
			return false
		}
	}
	return true
}

func EquivalentReadOnly(b1, b2 Board) bool {
	for i := range Size {
		readOnly := b1.IsReadOnly(i)
		if readOnly != b2.IsReadOnly(i) {
			return false
		}
		if readOnly && b1.Get(i) != b2.Get(i) {
			return false
		}
	}
	return true
}

// ContainsAll checks if all values in b2 are present in b1
func ContainsAll(b1, b2 Board) bool {
	for i := range Size {
		v := b2.Get(i)
		if v != 0 && b1.Get(i) != v {
			return false
		}
	}
	return true
}

// ContainsReadOnly checks if all read-only values in b2 are present in b1,
// ignoring edited values on both sides
func ContainsReadOnly(b1, b2 Board) bool {
	for i := range Size {
		if b2.IsReadOnly(i) && b1.Get(i) != b2.Get(i) {
			return false
		}
	}
	return true
}
