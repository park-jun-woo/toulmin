//ff:type feature=tangl type=model
//ff:what Node — a case node registration (`` `x` is a ... rule ... ``)
package ast

// Node is a single node registered inside a Case.
type Node struct {
	Name      string   `json:"name"`
	Role      Role     `json:"role"`
	Using     *Ref     `json:"using,omitempty"`    // nil implies a same-package function of the same name
	Checking  string   `json:"checking,omitempty"` // target case name for verdict composition
	With      []string `json:"with,omitempty"`     // Definitions term names
	Qualified *float64 `json:"qualified,omitempty"`
	Line      int      `json:"line"`
}
