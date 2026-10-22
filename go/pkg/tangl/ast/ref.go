//ff:type feature=tangl type=model
//ff:what Ref — reference to a name, optionally qualified by a package alias
package ast

// Ref is a reference to a name, optionally qualified by a package alias.
// The source `arm`.`placeCup` becomes {Alias:"arm", Name:"placeCup"};
// a local `fn` becomes {Name:"fn"}.
type Ref struct {
	Alias string `json:"alias,omitempty"`
	Name  string `json:"name"`
}
