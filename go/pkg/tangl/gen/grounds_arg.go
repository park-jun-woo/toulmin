//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what groundsArg — renders a Grounds name list as comma-separated Go string literals
package gen

import (
	"fmt"
	"strings"
)

// groundsArg renders the Grounds name list as the argument list passed to
// .Grounds(...): each name becomes a quoted Go string literal in document
// order (no lookup, no sorting, no dedup — see Phase016 determinism note).
func groundsArg(names []string) string {
	quoted := make([]string, len(names))
	for i, name := range names {
		quoted[i] = fmt.Sprintf("%q", name)
	}
	return strings.Join(quoted, ", ")
}
