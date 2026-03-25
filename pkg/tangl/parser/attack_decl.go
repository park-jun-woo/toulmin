//ff:type feature=tangl type=model
//ff:what AttackDecl — attack relationship declaration
package parser

// AttackDecl represents an attack: attacker attacks target.
type AttackDecl struct {
	Attacker string
	Target   string
	Graph    string // inferred from parent in nested lists
	Line     int
}
