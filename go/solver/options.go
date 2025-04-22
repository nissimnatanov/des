package solver

// Options for the user to start the solver.
type Options struct {
	Action Action
	// MaxRecursionDepth is the maximum recursion depth for the solver, defaults to 10
	MaxRecursionDepth int
}
