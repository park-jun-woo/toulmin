//ff:type feature=tangl type=model
//ff:what Import — package import declaration
package parser

// Import represents a package import: alias is from "path".
type Import struct {
	Alias   string
	Package string
	Line    int
}
