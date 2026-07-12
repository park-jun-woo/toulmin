//ff:func feature=tangl type=codegen control=sequence
//ff:what callFnSelector — renders the Go symbol a Rules host Call binds to
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// callFnSelector renders the Go symbol a Rules host Call binds to: an aliased
// reference stays "alias.Name" verbatim (a See-imported package symbol),
// while a bare reference becomes the goIdent local identifier — resolved by
// the Go linker to either the same-named generated tangl:Rules function or a
// hand-written package-level function (the document-local-first two-stage
// lookup, where a collision is a Go duplicate-identifier compile error).
func callFnSelector(c ast.Call) string {
	if c.Fn.Alias != "" {
		return c.Fn.Alias + "." + c.Fn.Name
	}
	return goIdent(c.Fn.Name)
}
