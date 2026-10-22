//ff:type feature=tangl type=model
//ff:what SpecRef — a reference to a with-injected spec argument, optionally path-extended
package ast

// SpecRef is a reference to a with-injected spec argument. The first step
// names the spec term; any remaining steps walk into it.
type SpecRef struct {
	Steps []string `json:"steps"`
}
