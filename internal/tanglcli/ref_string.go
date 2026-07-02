//ff:func feature=cli type=command control=sequence
//ff:what refString — renders a Ref as "alias.name" or bare "name"
package tanglcli

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// refString renders ref as "alias.name" when it has a package alias, or
// bare "name" for a local reference.
func refString(ref ast.Ref) string {
	if ref.Alias == "" {
		return ref.Name
	}
	return ref.Alias + "." + ref.Name
}
