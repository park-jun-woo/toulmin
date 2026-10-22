//ff:type feature=tangl type=model
//ff:what Require — a required context field declaration (`<field> is required`)
package ast

// Require is a `field` is required [as Type] declaration.
type Require struct {
	Field string `json:"field"`
	Type  string `json:"type,omitempty"`
	Line  int    `json:"line"`
}
