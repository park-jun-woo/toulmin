//ff:type feature=engine type=model
//ff:what EvalOption — evaluation method option (Matrix or Recursive)
package toulmin

// EvalOption controls how verdict is computed.
type EvalOption int

const (
	// Matrix uses matrix multiplication for verdict computation (default).
	Matrix EvalOption = iota
	// Recursive uses recursive h-Categoriser traversal for verdict computation.
	Recursive
)
