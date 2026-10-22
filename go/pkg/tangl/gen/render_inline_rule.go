//ff:func feature=tangl type=codegen control=sequence
//ff:what renderInlineRule — writes one tangl:Rules entry as a rule function
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderInlineRule writes one tangl:Rules entry as a
// func(toulmin.Context, toulmin.Specs) (bool, any) — the exact signature
// toulmin.Rule/Counter/Except require — whose body invokes the closure
// composition renderClosureExpr builds from the entry's condition tree,
// returning the verdict together with the auto-generated witness Fact of the
// short-circuit stopping leaf.
func renderInlineRule(w *strings.Builder, r ast.InlineRule) error {
	expr, err := renderClosureExpr(r.Cond, nil, 0)
	if err != nil {
		return fmt.Errorf("gen: rule %q: %w", r.Name, err)
	}
	fmt.Fprintf(w, "func %s(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {\n", goIdent(r.Name))
	fmt.Fprintf(w, "\treturn %s()\n", expr)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	return nil
}
