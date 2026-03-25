//ff:type feature=tangl type=model
//ff:what File — parsed TANGL file AST
package parser

// File represents a parsed TANGL markdown file.
type File struct {
	Imports  []Import
	Rules    []InlineRule
	Graphs   []GraphDecl
	Bindings []RuleBinding
	Attacks  []AttackDecl
	Evals    []EvalDecl
}
