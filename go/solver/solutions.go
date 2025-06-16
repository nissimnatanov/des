package solver

import (
	"slices"

	"github.com/nissimnatanov/des/go/boards"
)

type Solutions []*boards.Solution

func (s Solutions) Append(newSol *boards.Solution) Solutions {
	for _, sol := range s {
		if boards.Equivalent(sol, newSol) {
			return s
		}
	}
	return append(s, newSol)
}

func (s Solutions) With(other Solutions) Solutions {
	if len(s) == 0 {
		// assuming other is already unique
		return slices.Clone(other)
	}
	for _, sol := range other {
		s = s.Append(sol)
	}
	return s
}
