//ff:func feature=tangl type=validator control=sequence
//ff:what ruleCallAliasErrs — collect undeclared-alias errors for host calls inside one rule's condition
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// ruleCallAliasErrs walks one rule's condition tree and reports every host Call
// whose FUNC_REF carries a package alias not present in aliases. It covers both
// predicate-position Calls and value-position Calls (Compare operands and an
// Exists source).
func ruleCallAliasErrs(doc *ast.Document, r ast.InlineRule, aliases map[string]bool) []error {
	var errs []error
	report := func(fn ast.Ref) {
		if fn.Alias != "" && !aliases[fn.Alias] {
			errs = append(errs, errAt(doc.Path, r.Line, "rule %q calls host function via undeclared package alias %q", r.Name, fn.Alias))
		}
	}
	reportOperand := func(o ast.Operand) {
		walkOperand(o, func(op ast.Operand) {
			if c, ok := op.(ast.Call); ok {
				report(c.Fn)
			}
		})
	}
	walkExpr(r.Cond, func(e ast.Expr) {
		switch v := e.(type) {
		case ast.Call:
			report(v.Fn)
		case ast.Compare:
			reportOperand(v.Lhs)
			reportOperand(v.Rhs)
		case ast.Exists:
			reportOperand(v.Source)
		}
	})
	return errs
}
