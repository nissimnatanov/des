package solver

import (
	"time"

	"github.com/nissimnatanov/des/go/boards"
)

type Result struct {
	// Input is the initial input board, as provided by the caller.
	Input  *boards.Game `json:"board,omitempty"`
	Action Action       `json:"action,omitempty"`
	// Play is the player's board used by the solver, it includes the progress made on the board
	// during the solving process. If the solving process was successful, this board will be
	// equivalent to the solution. If the solving process failed, this board will contain the
	// progress made until the failure point.
	Play      *boards.Game  `json:"-"`
	Status    Status        `json:"status"`
	Error     error         `json:"error,omitempty"`
	Steps     StepStats     `json:"steps"`
	Solutions *Solutions    `json:"-"`
	Elapsed   time.Duration `json:"elapsed"`
}

func (r *Result) complete(status Status) *Result {
	r.Status = status
	return r
}

func (r *Result) completeErr(err error) *Result {
	r.Status = StatusError
	r.Error = err
	return r
}
