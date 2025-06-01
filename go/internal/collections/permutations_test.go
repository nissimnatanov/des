package collections_test

import (
	"slices"
	"testing"

	"github.com/nissimnatanov/des/go/internal/collections"
	"gotest.tools/v3/assert"
)

func collect[T any](items []T) [][]T {
	return slices.Collect(collections.MapIter(collections.Permutate(items), slices.Clone))
}

func TestPermutations(t *testing.T) {
	assert.DeepEqual(t, collect([]int(nil)), [][]int(nil))
	assert.DeepEqual(t, collect([]int{}), [][]int(nil))

	assert.DeepEqual(t, collect([]int{22}), [][]int{{22}})
	assert.DeepEqual(t, collect([]int{22, 33}), [][]int{{22, 33}, {33, 22}})

	assert.DeepEqual(t, collect([]int{1, 2, 3}),
		[][]int{
			{1, 2, 3}, {1, 3, 2},
			{2, 1, 3}, {2, 3, 1},
			{3, 2, 1}, {3, 1, 2},
		})

	assert.DeepEqual(t, collect([]string{"1", "2", "3", "4"}),
		[][]string{
			{"1", "2", "3", "4"}, {"1", "2", "4", "3"},
			{"1", "3", "2", "4"}, {"1", "3", "4", "2"},
			{"1", "4", "3", "2"}, {"1", "4", "2", "3"},
			{"2", "1", "3", "4"}, {"2", "1", "4", "3"},
			{"2", "3", "1", "4"}, {"2", "3", "4", "1"},
			{"2", "4", "3", "1"}, {"2", "4", "1", "3"},
			{"3", "2", "1", "4"}, {"3", "2", "4", "1"},
			{"3", "1", "2", "4"}, {"3", "1", "4", "2"},
			{"3", "4", "1", "2"}, {"3", "4", "2", "1"},
			{"4", "2", "3", "1"}, {"4", "2", "1", "3"},
			{"4", "3", "2", "1"}, {"4", "3", "1", "2"},
			{"4", "1", "3", "2"}, {"4", "1", "2", "3"},
		})
}
