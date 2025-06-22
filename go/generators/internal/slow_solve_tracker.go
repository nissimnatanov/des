package internal

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nissimnatanov/des/go/boards"
	"github.com/nissimnatanov/des/go/solver"
)

const MaxSlowSolveBoards = 5

var SlowBoards slowSolveTracker

type slowSolveTracker struct {
	slow []*solver.Result
}

func (s *slowSolveTracker) Add(res *solver.Result) {
	if len(s.slow) < MaxSlowSolveBoards {
		s.slow = append(s.slow, res)
		s.sort()
	} else if s.slow[len(s.slow)-1].Elapsed < res.Elapsed {
		s.slow[len(s.slow)-1] = res
		s.sort()
	}
}

func (s *slowSolveTracker) sort() {
	slices.SortFunc(s.slow, func(a, b *solver.Result) int {
		return int(b.Elapsed - a.Elapsed)
	})
}

func (s *slowSolveTracker) Log() string {
	logs := make([]string, 0, len(s.slow))
	for _, res := range s.slow {
		logs = append(logs,
			fmt.Sprintf("{ Board: %s, Elapsed: %s, Complexity: %d, Status: %s }",
				boards.Serialize(res.Input), res.Elapsed.Round(time.Millisecond/10), res.Complexity, res.Status))
	}
	return fmt.Sprintf("Slow boards: {\n  %s\n}", strings.Join(logs, "\n  "))
}
