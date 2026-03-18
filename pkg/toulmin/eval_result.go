//ff:type feature=engine type=model
//ff:what EvalResult — verdict result for a warrant node
package toulmin

// EvalResult holds the verdict and trace for an evaluated warrant.
type EvalResult struct {
	Name    string       `json:"name"`
	Verdict float64      `json:"verdict"`
	Trace   []TraceEntry `json:"trace"`
}
