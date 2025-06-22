package internal

import "github.com/nissimnatanov/des/go/solver"

func MaxComplexity(bestComplexity solver.StepComplexity, bs *BoardState) solver.StepComplexity {
	if bs == nil {
		return bestComplexity
	}
	return max(bestComplexity, bs.Complexity())
}

func MaxComplexityWithSorted(bestComplexity solver.StepComplexity, sbs *SortedBoardStates) solver.StepComplexity {
	if sbs == nil || sbs.Size() == 0 {
		return bestComplexity
	}
	return MaxComplexity(bestComplexity, sbs.Get(0))
}
