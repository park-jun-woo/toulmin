//ff:type feature=tangl type=model
//ff:what Endpoint — a tangl:Provides entry (external entry point)
package ast

// Endpoint is a single "provides `name`" entry: an external entry point
// that runs or checks Cases.
type Endpoint struct {
	Name     string    `json:"name"`
	Requires []Require `json:"requires,omitempty"`
	Runs     []string  `json:"runs,omitempty"`   // Run-mode case names
	Checks   []string  `json:"checks,omitempty"` // Evaluate-mode case names
	Line     int       `json:"line"`
}
