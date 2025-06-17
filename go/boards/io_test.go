package boards_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"gotest.tools/v3/assert"
)

// 3032828	       391.5 ns/op	     152 B/op	       5 allocs/op
// 1885412	       631.8 ns/op	     104 B/op	       3 allocs/op

func BenchmarkSerialize(b *testing.B) {
	in := "1D7A9B3B2C8B96B5D53B9C1B8C26D4C3F1B4F7B7C3B"
	board, err := boards.Deserialize(in)
	assert.NilError(b, err)
	bs := boards.Serialize(board)
	assert.Equal(b, bs, in)

	board.DisallowValue(1, 1)
	board.DisallowValue(80, 3)
	bs = boards.Serialize(board)
	assert.Equal(b, bs, "1[1]C7A9B3B2C8B96B5D53B9C1B8C26D4C3F1B4F7B7C3A[3]")

	for b.Loop() {
		boards.Serialize(board)
	}
}
