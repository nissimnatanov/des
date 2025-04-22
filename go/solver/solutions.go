package solver

import "github.com/nissimnatanov/des/go/board"

type Solutions struct {
	all []board.Solution
}

func (s *Solutions) Add(newSol board.Solution) {
	for _, sol := range s.all {
		if board.Equivalent(sol, newSol) {
			return
		}
	}
	s.all = append(s.all, newSol)
}

func (s *Solutions) Size() int {
	return len(s.all)
}

func (s *Solutions) At(i int) board.Solution {
	return s.all[i]
}
