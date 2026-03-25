//ff:type feature=tangl type=model
//ff:what GraphDecl — graph declaration
package parser

// GraphDecl represents a graph declaration: name is a graph "id".
type GraphDecl struct {
	Name string
	ID   string
	Line int
}
