//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what buildRules — renders every tangl:Rules entry in document order
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildRules renders every tangl:Rules entry as a standalone rule
// function matching the toulmin.Rule fn signature, in document order.
func buildRules(w *strings.Builder, rules []ast.InlineRule) error {
	for _, r := range rules {
		if err := renderInlineRule(w, r); err != nil {
			return err
		}
	}
	return nil
}
