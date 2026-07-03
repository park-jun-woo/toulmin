//ff:func feature=tangl type=parser control=selection pattern=recursive
//ff:what renderTerm — render one condition-list term (compare / not / either) as list-marker lines
package en

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderTerm renders a single condition-list term at depth, prefixed by its
// chaining word ("and"/"or", or none). A Compare or an inline-negated Compare is
// one line; an either group and a block "not" emit a head line plus their inner
// condition list at depth+1. The 2-space indentation and `- ` marker are relative;
// the core renderer shifts the whole subtree to the section base.
func renderTerm(prefix string, term ast.Expr, depth int) []string {
	head := strings.Repeat("  ", depth) + "- "
	if prefix != "" {
		head += prefix + " "
	}
	switch t := term.(type) {
	case ast.Compare:
		return []string{head + renderCompare(t)}
	case ast.Not:
		if inner, ok := t.Term.(ast.Compare); ok {
			return []string{head + "not " + renderCompare(inner)}
		}
		return append([]string{head + "not"}, renderCondList(t.Term, depth+1)...)
	case ast.Either:
		return append([]string{head + "either"}, renderCondList(t.Terms[0], depth+1)...)
	default:
		return []string{head + "`?`"}
	}
}
