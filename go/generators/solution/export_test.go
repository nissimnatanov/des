package solution

import (
	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/internal/random"
)

func GenerateSolutionWithCustomOrder(r *random.Random, sqOrder []int) *boards.Solution {
	if r == nil {
		r = random.New()
	}
	g := solutionGenerator{rand: r}
	return g.generate(sqOrder)
}
