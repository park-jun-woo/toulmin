//ff:type feature=tangl type=model
//ff:what SpecCall — spec factory call in a with clause
package parser

// SpecCall represents a spec factory call: Role("admin"), policy.IPList("block").
type SpecCall struct {
	Name string // Role, policy.IPList
	Args []any  // "admin", "block", 0.5, etc.
}
