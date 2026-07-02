//ff:func feature=tangl type=codegen control=sequence
//ff:what renderInlineRule — writes one tangl:Rules entry as a rule function
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderInlineRule writes one tangl:Rules entry as a
// func(toulmin.Context, toulmin.Specs) (bool, any) — the exact signature
// toulmin.Rule/Counter/Except require — whose body evaluates the entry's
// condition tree via renderExpr.
func renderInlineRule(w *strings.Builder, r ast.InlineRule) error {
	cond, err := renderExpr(r.Cond)
	if err != nil {
		return fmt.Errorf("gen: rule %q: %w", r.Name, err)
	}
	fmt.Fprintf(w, "func %s(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {\n", goIdent(r.Name))
	fmt.Fprintf(w, "\treturn %s, nil\n", cond)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	return nil
}
