//ff:type feature=engine type=model
//ff:what TraceEntry — single rule evaluation record in trace
package toulmin

import "time"

// TraceEntry records one rule's evaluation result for explainability.
type TraceEntry struct {
	Name      string        `json:"name"`
	Role      string        `json:"role"`
	Activated bool          `json:"activated"`
	Qualifier float64       `json:"qualifier"`
	Evidence  any           `json:"evidence,omitempty"`
	Backing   any           `json:"backing,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
}
