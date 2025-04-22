package solver

type Result struct {
	Status    Status
	Error     error
	StepStats StepStats
	Solutions *Solutions
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
