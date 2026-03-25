//ff:type feature=tangl type=model
//ff:what InlineRule — inline rule definition with expression
package parser

// InlineRule represents an inline rule: rule "name" is return that ...
type InlineRule struct {
	Name string
	Expr Expr
	Line int
}
