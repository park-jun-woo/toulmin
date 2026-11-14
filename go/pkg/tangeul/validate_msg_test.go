//ff:func feature=tangl type=util control=sequence
//ff:what validateMsg — normalizeLoc'd validate.Validate message ("" when the AST is clean)
package tangeul

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/validate"
)

// validateMsg runs validate.Validate and returns its error text with every
// "path:line:" prefix rewritten by normalizeLoc, or "" when the document is clean.
// The normalization lets the two routes (which embed different paths and line
// numbers) compare on semantic content alone.
func validateMsg(a *ast.Document) string {
	err := validate.Validate(a)
	if err == nil {
		return ""
	}
	return normalizeLoc(err.Error())
}
