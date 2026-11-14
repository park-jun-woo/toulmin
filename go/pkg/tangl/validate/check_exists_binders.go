//ff:func feature=tangl type=validator control=iteration dimension=1 pattern=recursive
//ff:what checkExistsBinders — rejects empty exists binders and binders shadowing an enclosing binder
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkExistsBinders reports an error for every Exists whose binder name is
// empty, or whose binder shadows an enclosing binder of the same name in a
// nested exists. Shadowing is an error (not a lint warning) because a silently
// captured outer element would change the judgment's meaning.
func checkExistsBinders(doc *ast.Document) []error {
	var errs []error
	var visit func(e ast.Expr, active map[string]bool, rule string, line int)
	visit = func(e ast.Expr, active map[string]bool, rule string, line int) {
		switch v := e.(type) {
		case ast.Exists:
			if v.Binder == "" {
				errs = append(errs, errAt(doc.Path, line, "rule %q: exists quantifier has an empty binder name", rule))
			} else if active[v.Binder] {
				errs = append(errs, errAt(doc.Path, line, "rule %q: exists binder %q shadows an enclosing binder of the same name", rule, v.Binder))
			}
			inner := make(map[string]bool, len(active)+1)
			for k := range active {
				inner[k] = true
			}
			inner[v.Binder] = true
			visit(v.Cond, inner, rule, line)
		case ast.Logic:
			for _, t := range v.Terms {
				visit(t, active, rule, line)
			}
		case ast.Either:
			for _, t := range v.Terms {
				visit(t, active, rule, line)
			}
		case ast.Not:
			visit(v.Term, active, rule, line)
		}
	}
	for _, r := range doc.Rules {
		visit(r.Cond, map[string]bool{}, r.Name, r.Line)
	}
	return errs
}
