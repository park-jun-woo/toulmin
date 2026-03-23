//ff:func feature=analyzer type=analyzer control=sequence
//ff:what collectRuleName — extracts rule function name and role from Warrant/Rebuttal/Defeater call
package analyzer

import (
	"go/ast"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// collectRuleName extracts the first argument's identifier name and role from a rule registration call.
func collectRuleName(call *ast.CallExpr, gc *graphCollector, method string) {
	if len(call.Args) == 0 {
		return
	}
	ident, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return
	}
	gc.def.Rules = append(gc.def.Rules, toulmin.GraphRuleDef{
		Name:      ident.Name,
		Role:      strings.ToLower(method),
		Qualifier: 1.0,
	})
}
