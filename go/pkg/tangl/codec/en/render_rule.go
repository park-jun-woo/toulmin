//ff:func feature=tangl type=parser control=selection
//ff:what renderRule — render an inline Rules entry (one-line lone compare/exists/call or multi-line nested tree)
package en

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderRule renders a tangl:Rules entry. A lone Compare, a lone Exists with a
// lone-compare body, or a lone predicate Call round-trips one-line ("`name` when
// <compare>" / "`name` when any of ... satisfies <compare>" / "`name` when `f`
// holds"); and/or chains, either, not, and nested-body exists use the multi-line
// head + nested list. This mirrors the parser's rule-head dispatch (inline_term =
// exists_term | call_term | compare), so the canonical form re-parses to the same
// tree from its minimal shape.
func renderRule(r ast.InlineRule) string {
	if r.Using != nil {
		s := "`" + r.Name + "` uses " + renderRef(*r.Using)
		s += renderWith(r.With)
		if len(r.Grounds) > 0 {
			s += " grounded on " + renderBacktickNames(r.Grounds)
		}
		return s
	}
	head := "`" + r.Name + "` when"
	switch c := r.Cond.(type) {
	case ast.Compare:
		return head + " " + renderCompare(c)
	case ast.Exists:
		return strings.Join(renderExists(head+" ", c, 0), "\n")
	case ast.Call:
		return head + " " + renderRef(c.Fn) + " holds"
	}
	lines := renderCondList(r.Cond, 1)
	return head + "\n" + strings.Join(lines, "\n")
}

func renderBacktickNames(names []string) string {
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = "`" + name + "`"
	}
	return strings.Join(parts, " and ")
}
