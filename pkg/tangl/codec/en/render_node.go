//ff:func feature=tangl type=parser control=selection
//ff:what renderNode — render a case node registration (`` `X` is a/an <role> rule [using|checking ...]``)
package en

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderNode renders a case node back to "`<name>` is <a|an> <role> rule" plus its
// optional trailing clause: "using <ref>" with its "with `t` and `u`" and
// "qualified <float>" sub-clauses, or "checking `<case>`". The article is "an" for
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
		if n.Qualified != nil {
			s += " qualified " + strconv.FormatFloat(*n.Qualified, 'g', -1, 64)
		}
	} else if n.Checking != "" {
		s += " checking `" + n.Checking + "`"
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
