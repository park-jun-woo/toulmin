//ff:type feature=engine type=model
//ff:what EvalMethod — verdict computation method (Matrix or Recursive)
package toulmin

// EvalMethod controls which algorithm computes the verdict.
type EvalMethod int

const (
	// Matrix uses matrix multiplication for verdict computation (default).
	Matrix EvalMethod = iota
	// Recursive uses recursive h-Categoriser traversal for verdict computation.
	Recursive
)
