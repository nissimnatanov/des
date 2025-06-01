package collections

import (
	"slices"
	"strconv"
)

func permuteRecursive[T any, S ~[]T](items []T, start int, yield func(S) bool) bool {
	if start >= len(items) {
		// Send the current permutation on the channel
		return yield(items)
	}
	for i := start; i < len(items); i++ {
		items[start], items[i] = items[i], items[start]
		if !permuteRecursive(items, start+1, yield) {
			// stop if the yield function returns false
			return false
		}
		items[start], items[i] = items[i], items[start] // Backtrack
	}
	return true
}

// ClonePermutations generates all permutations of the given slice of items, cloning
// the slice each time it is yielded. This is useful when the yield func needs to
// operate on the clone.
func ClonePermutations[T any, S ~[]T](items S) func(yield func(S) bool) {
	return MapIter(Permutate(slices.Clone(items)), slices.Clone)
}

// Permutate generates all permutations of the given slice of items, inline,
// yield slice is always called on the same original slice. If the yield function
// needs to keep the current permutation, it must clone it before return.
//
// If the original slice must be preserved, pass a clone. Also, if yield must receive
// a clone, consider using ClonePermutations instead.
func Permutate[T any, S ~[]T](items S) func(yield func(S) bool) {
	return func(yield func(S) bool) {
		switch len(items) {
		case 0:
			return
		case 1:
			yield(items)
			return
		case 2:
			if yield(items) {
				items[0], items[1] = items[1], items[0]
				yield(items)
			}
			return
		}
		const recursionCap = 20
		if len(items) > recursionCap {
			// do not want to cause a stack overflow, so for now just hard-cap it
			// if needed, implement a non-recursive version here
			panic("permutationsOf: too many items, max is " + strconv.Itoa(recursionCap))
		}
		permuteRecursive(items, 0, yield)
	}
}
