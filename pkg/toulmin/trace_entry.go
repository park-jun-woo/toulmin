//ff:type feature=engine type=model
//ff:what TraceEntry — single rule evaluation record (Name=Claim, Ground=ctx, Specs=Backing)
package toulmin

import "time"

// TraceEntry records one rule's evaluation result for explainability.
// Name is the Claim, Ground is the ctx the rule saw, Specs is the Backing.
type TraceEntry struct {
	Name      string        `json:"name"`             // = Claim
	Role      string        `json:"role"`
	Activated bool          `json:"activated"`
	Qualifier float64       `json:"qualifier"`
	Verdict   float64       `json:"verdict"`
	Evidence  any           `json:"evidence,omitempty"`
	Ground    any           `json:"ground,omitempty"` // = Ground (ctx as-is)
	Specs     Specs         `json:"specs,omitempty"`  // = Backing
	Duration  time.Duration `json:"duration,omitempty"`
}
