//ff:type feature=tangl type=model
//ff:what Supersede — a supersede edge ("`x` supersedes `y`")
package ast

// Supersede is a supersede edge: superseder.Supersedes(superseded).
// Unlike Attack it crisply excludes the target from the tally rather than
// lowering its verdict.
type Supersede struct {
	Superseder string `json:"superseder"`
	Superseded string `json:"superseded"`
	Line       int    `json:"line"`
}
