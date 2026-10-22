//ff:func feature=tangl type=util control=sequence
//ff:what callRepr — renders a Call reference as "alias.name" or "name"
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// callRepr renders a Call's reference: "alias.name" when the reference is
// alias-qualified, otherwise the bare "name".
func callRepr(c ast.Call) string {
	if c.Fn.Alias != "" {
		return c.Fn.Alias + "." + c.Fn.Name
	}
	return c.Fn.Name
}
