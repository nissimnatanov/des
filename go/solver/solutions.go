package solver

import "github.com/nissimnatanov/des/go/boards"

type Solutions struct {
	all []*boards.Solution
}

func (s *Solutions) Add(newSol *boards.Solution) {
	for _, sol := range s.all {
		if boards.Equivalent(sol, newSol) {
			return
		}
	}
	s.all = append(s.all, newSol)
}

func (s *Solutions) Size() int {
	return len(s.all)
}

func (s *Solutions) At(i int) *boards.Solution {
	return s.all[i]
}
