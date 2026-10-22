//ff:type feature=tangl type=model
//ff:what InlineRule — a tangl:Rules entry (`` `name` when <condition> ``)
package ast

// InlineRule is a single entry of the tangl:Rules section.
type InlineRule struct {
	Name    string   `json:"name"`
	Cond    Expr     `json:"cond,omitempty"`
	Using   *Ref     `json:"using,omitempty"`
	With    []string `json:"with,omitempty"`
	Grounds []string `json:"grounds,omitempty"`
	Line    int      `json:"line"`
}
