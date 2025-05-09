package solver

type Result struct {
	Status    Status     `json:"status"`
	Error     error      `json:"error,omitempty"`
	Steps     StepStats  `json:"steps"`
	Solutions *Solutions `json:"-"`
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
