//ff:type feature=scanner type=model
//ff:what RuleDecl — extracted rule declaration (func name + annotation lines)
package scanner

// RuleDecl holds a function name and its //rule: annotation lines.
type RuleDecl struct {
	FuncName string
	Lines    []string
}
