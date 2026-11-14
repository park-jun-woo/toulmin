//ff:func feature=tangl type=engine control=sequence
//ff:what pathRepr — renders a Path as dotted steps, or as source[index].rest under a binder
package eval

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// pathRepr renders a Path's Where string. When the first step names an
// in-scope exists binder it renders as "<source>[<index>].<rest>"; otherwise it
// is the dotted join of the steps.
func pathRepr(p ast.Path, sc *scope) string {
	if len(p.Steps) == 0 {
		return ""
	}
	b, ok := sc.binders[p.Steps[0]]
	if !ok {
		return strings.Join(p.Steps, ".")
	}
	base := reprOperand(b.source, sc) + "[" + strconv.Itoa(b.index) + "]"
	if len(p.Steps) == 1 {
		return base
	}
	return base + "." + strings.Join(p.Steps[1:], ".")
}
