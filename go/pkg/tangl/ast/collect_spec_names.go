//ff:func feature=tangl type=model control=selection pattern=recursive
//ff:what collectSpecNames — walk a condition tree appending unseen first-step SpecRef names
package ast

// collectSpecNames visits expr in source order and appends the first step of
// every SpecRef operand not yet in seen. Compare checks Lhs then Rhs; Logic and
// Either recurse into their Terms; Not recurses into its Term; Exists records
// its Source before recursing into its Cond.
func collectSpecNames(expr Expr, seen map[string]bool, names []string) []string {
	switch value := expr.(type) {
	case Compare:
		if value.Lhs != nil {
			names = addSpecOperand(value.Lhs, seen, names)
		}
		if value.Rhs != nil {
			names = addSpecOperand(value.Rhs, seen, names)
		}
	case Logic:
		names = collectSpecNameTerms(value.Terms, seen, names)
	case Either:
		names = collectSpecNameTerms(value.Terms, seen, names)
	case Not:
		names = collectSpecNames(value.Term, seen, names)
	case Exists:
		names = addSpecOperand(value.Source, seen, names)
		names = collectSpecNames(value.Cond, seen, names)
	}
	return names
}
