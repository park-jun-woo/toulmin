//ff:func feature=tangl type=util control=sequence
//ff:what refString — formats an ast.Ref as "alias.name" or "name"
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// refString formats r as "alias.name" when aliased, or as "name" for a
// local (unaliased) reference.
func refString(r *ast.Ref) string {
	if r == nil {
		return ""
	}
	if r.Alias == "" {
		return r.Name
	}
	return r.Alias + "." + r.Name
}
