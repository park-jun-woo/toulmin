//ff:type feature=engine type=model
//ff:what EvalResult — verdict result for a warrant node
package toulmin

// EvalResult holds the verdict and trace for an evaluated warrant.
type EvalResult struct {
	Name     string  `json:"name"`
	Verdict  float64 `json:"verdict"`
	Evidence any     `json:"evidence,omitempty"`
	// Applicability is the warrant's effective applicability a_eff in [0, 1]:
	// (v+1)/2 of its own rule verdict, multiplied by every Except gate (1 − F·(v_j+1)/2).
	Applicability float64      `json:"applicability"`
	Trace         []TraceEntry `json:"trace"`
}
