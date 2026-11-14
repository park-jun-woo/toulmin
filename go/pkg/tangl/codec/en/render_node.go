//ff:func feature=tangl type=parser control=selection
//ff:what renderNode — render a case node registration (`` `X` is a/an <role> rule [using|checking ...] [qualified F] [, level ...]``)
package en

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderNode renders a case node back to "`<name>` is <a|an> <role> rule" plus its
// optional trailing clauses: "using <ref>" with its "with `t` and `u`" sub-clause,
// or "checking <case-ref>" (alias-qualified when Checking.Alias is set), then
// "qualified <float>" (also without using/checking — the canonical binding of a
// Rules-declared predicate), the legacy using form's "grounded on" list, then
// ", level fail|review". The article is "an" for
// an except rule and "a" otherwise. It inverts parseNodeItem/applyNodeClause.
func renderNode(n ast.Node) string {
	article, role := "a", "general rule"
	switch n.Role {
	case ast.CounterRule:
		role = "counter rule"
	case ast.ExceptRule:
		article, role = "an", "except rule"
	}
	s := "`" + n.Name + "` is " + article + " " + role
	if n.Using != nil {
		s += " using " + renderRef(*n.Using) + renderWith(n.With)
	} else if n.Checking != nil {
		s += " checking " + renderRef(*n.Checking)
	}
	if n.Qualified != nil {
		s += " qualified " + strconv.FormatFloat(*n.Qualified, 'g', -1, 64)
	}
	if n.Using != nil && len(n.Grounds) > 0 {
		s += " grounded on " + renderBacktickNames(n.Grounds)
	}
	if n.Level != ast.LevelUnset {
		word := "fail"
		if n.Level == ast.LevelReview {
			word = "review"
		}
		s += ", level " + word
	}
	return s
}
