//ff:type feature=engine type=model
//ff:what EvalMethod — verdict computation method (Matrix or Recursive)
package toulmin

// EvalMethod controls which algorithm computes the verdict.
type EvalMethod int

const (
	// Matrix uses lazy recursive h-Categoriser for verdict computation (default).
	Matrix EvalMethod = iota
	// Recursive is not yet implemented.
	Recursive
)
