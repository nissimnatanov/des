package boards_test

import (
	"testing"

	"github.com/nissimnatanov/des/go/boards"
	"gotest.tools/v3/assert"
)

func TestDeserializeSolution(t *testing.T) {
	sol, err := boards.DeserializeSolution("534678912 672195348 198342567 859761423 426853791 713924856 961537284 287419635 345286179")
	assert.NilError(t, err)

	b := sol.Clone(boards.EditMode)
	assert.Assert(t, b != nil)
	assert.Assert(t, b.IsSolved(), "board is not solved: %q", b)
	// ╔═══════╦═══════╦═══════╗
	// ║ 5.3.4.║ 6.7.8.║ 9.1.2.║
	// ║ 6.7.2.║ 1.9.5.║ 3.4.8.║
	// ║ 1.9.8.║ 3.4.2.║ 5.6.7.║
	// ╠═══════╬═══════╬═══════╣
	// ║ 8.5.9.║ 7.6.1.║ 4.2.3.║
	// ║ 4.2.6.║ 8.5.3.║ 7.9.1.║
	// ║ 7.1.3.║ 9.2.4.║ 8.5.6.║
	// ╠═══════╬═══════╬═══════╣
	// ║ 9.6.1.║ 5.3.7.║ 2.8.4.║
	// ║ 2.8.7.║ 4.1.9.║ 6.3.5.║
	// ║ 3.4.5.║ 2.8.6.║ 1.7.9.║
	// ╚═══════╩═══════╩═══════╝

	b = sol.Clone(boards.PlayMode)
	assert.Assert(t, b != nil)
	assert.Assert(t, boards.ContainsAll(sol, b))
}
