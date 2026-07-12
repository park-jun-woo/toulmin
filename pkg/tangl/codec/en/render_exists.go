//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what renderExists — render an exists quantifier (one-line compare body, or head + nested list)
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// renderExists renders an ast.Exists back to its English surface form, shared by
// the term position (renderTerm) and the rule-head position (renderRule) so the
// two dispatch sites cannot drift. lead is the leading head text — the "- "
// list marker for a term, or "`name` when " for a rule head. When Cond is a lone
// Compare the body follows "satisfies" on the same line; otherwise the head line
// is followed by the nested condition list at depth+1.
func renderExists(lead string, t ast.Exists, depth int) []string {
	eh := lead + "any of " + renderOperand(t.Source) + " as `" + t.Binder + "` satisfies"
	if cmp, ok := t.Cond.(ast.Compare); ok {
		return []string{eh + " " + renderCompare(cmp)}
	}
	return append([]string{eh}, renderCondList(t.Cond, depth+1)...)
}
