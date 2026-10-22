//ff:func feature=tangl type=model control=iteration dimension=2 pattern=recursive
//ff:what RuleSpecNames — collect first-step SpecRef names from a rule condition
package ast

// RuleSpecNames returns the first step of every SpecRef in a condition, once,
// in source traversal order. These names form an expression rule's inferred
// definition bindings.
func RuleSpecNames(rule InlineRule) []string {
	var names []string
	seen := map[string]bool{}
	addOperand := func(op Operand) {
		if ref, ok := op.(SpecRef); ok && len(ref.Steps) > 0 && !seen[ref.Steps[0]] {
			seen[ref.Steps[0]] = true
			names = append(names, ref.Steps[0])
		}
	}
	var visit func(Expr)
	visit = func(expr Expr) {
		switch value := expr.(type) {
		case Compare:
			if value.Lhs != nil {
				addOperand(value.Lhs)
			}
			if value.Rhs != nil {
				addOperand(value.Rhs)
			}
		case Logic:
			for _, term := range value.Terms {
				visit(term)
			}
		case Either:
			for _, term := range value.Terms {
				visit(term)
			}
		case Not:
			visit(value.Term)
		case Exists:
			addOperand(value.Source)
			visit(value.Cond)
		}
	}
	visit(rule.Cond)
	return names
}
