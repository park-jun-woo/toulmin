//ff:type feature=tangl type=model
//ff:what Field — a struct definition field (`has `f` as Type`)
package ast

// Field is a single field of a struct Definition.
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Line int    `json:"line"`
}
