//ff:func feature=tangl type=parser control=selection
//ff:what appendNode — append a parsed AST node to the appropriate File field
package parser

// appendNode appends a parsed AST node to the appropriate slice on File.
func appendNode(f *File, node any) {
	switch v := node.(type) {
	case Import:
		f.Imports = append(f.Imports, v)
	case GraphDecl:
		f.Graphs = append(f.Graphs, v)
	case RuleBinding:
		f.Bindings = append(f.Bindings, v)
	case AttackDecl:
		f.Attacks = append(f.Attacks, v)
	case EvalDecl:
		f.Evals = append(f.Evals, v)
	case InlineRule:
		f.Rules = append(f.Rules, v)
	}
}
