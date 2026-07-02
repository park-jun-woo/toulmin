//ff:func feature=tangl type=codegen control=sequence
//ff:what resolveSeePath — maps a tangl:See package path to its Go import path
package gen

import "strings"

// resolveSeePath maps a tangl:See package path to its Go import path: a
// "tangl/..." path resolves inside this module ("tangl/log" becomes
// "github.com/park-jun-woo/toulmin/pkg/tangl/log"); any other path is
// used verbatim as an external import path.
func resolveSeePath(pkgPath string) string {
	const prefix = "tangl/"
	if strings.HasPrefix(pkgPath, prefix) {
		return "github.com/park-jun-woo/toulmin/pkg/tangl/" + strings.TrimPrefix(pkgPath, prefix)
	}
	return pkgPath
}
