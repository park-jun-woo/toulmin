//ff:type feature=tangl type=model
//ff:what Role — case node role (general/counter/except rule)
package ast

// Role classifies a Node's engine role.
type Role int

const (
	// GeneralRule maps to a toulmin Rule (Warrant).
	GeneralRule Role = iota
	// CounterRule maps to a toulmin Counter (Rebuttal).
	CounterRule
	// ExceptRule maps to a toulmin Except (Defeater).
	ExceptRule
)
