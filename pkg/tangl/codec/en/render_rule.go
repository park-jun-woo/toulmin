//ff:func feature=tangl type=parser control=sequence
//ff:what renderRule — render an inline Rules entry (one-line single compare or multi-line nested tree)
package en

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderRule renders a tangl:Rules entry. A single leaf Compare uses the one-line
// form "`name` when <compare>"; every other condition (a Logic and/or chain, an
// either group, or a not) uses the multi-line form — a "`name` when" head plus the
// nested condition list. Only a lone Compare round-trips one-line, since the
// parser's inline branch (parseCompare) parses exactly one comparison; and/or
// chains, either, and not must be nested list items to re-parse to the same tree.
func renderRule(r ast.InlineRule) string {
	head := "`" + r.Name + "` when"
	if cmp, ok := r.Cond.(ast.Compare); ok {
		return head + " " + renderCompare(cmp)
	}
	lines := renderCondList(r.Cond, 1)
	return head + "\n" + strings.Join(lines, "\n")
}
