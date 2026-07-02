//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeCaseDoc — writes a doc comment naming a case's required ctx fields
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeCaseDoc writes a doc comment above a case's builder function
// naming the required ctx fields declared by its "is required" clauses,
// if any (informational only — Endpoint.Requires drives the actual
// tangl.Required call).
func writeCaseDoc(w *strings.Builder, c ast.Case) {
	if len(c.Requires) == 0 {
		return
	}
	fields := make([]string, len(c.Requires))
	for i, r := range c.Requires {
		fields[i] = r.Field
	}
	fmt.Fprintf(w, "// requires ctx fields: %s\n", strings.Join(fields, ", "))
}
