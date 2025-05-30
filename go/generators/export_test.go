package generators

import (
	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/generators/internal"
)

const FastGenerationCap = fastGenerationCap

func GenerateSolutionWithCustomOrder(r *internal.Random, sqOrder []int) *boards.Solution {
	if r == nil {
		r = internal.NewRandom()
	}
	return solutionGenerator{rand: r}.generate(sqOrder)
}
