//ff:type feature=tangl type=model
//ff:what Path — a ctx data path: backtick names chained by dots (struct field / map key steps)
package ast

// Path is a ctx data path: backtick names chained by dots, each step a
// struct-field or map-key access (e.g. ["author", "trustScore"]). A leading
// exists binder step resolves to the bound element instead of a ctx key.
type Path struct {
	Steps []string `json:"steps"`
}
