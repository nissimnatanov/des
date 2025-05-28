package generators

import "github.com/nissimnatanov/des/go/boards"

func GenerateSolutionWithCustomOrder(r *Random, sqOrder []int) *boards.Solution {
	if r == nil {
		r = NewRandom()
	}
	return solutionGenerator{rand: r}.generate(sqOrder)
}
