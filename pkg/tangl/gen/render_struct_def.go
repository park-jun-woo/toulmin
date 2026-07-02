//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what renderStructDef — writes one StructDef as a Go struct type
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderStructDef writes one StructDef as a Go struct type: each "has"
// field becomes an exported (PascalCase) field, typed by goFieldType.
func renderStructDef(w *strings.Builder, d ast.Definition) defInfo {
	fmt.Fprintf(w, "type %s struct {\n", goIdentExported(d.Name))
	for _, f := range d.Fields {
		fmt.Fprintf(w, "\t%s %s\n", goIdentExported(f.Name), goFieldType(f.Type))
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	return defInfo{Kind: ast.StructDef}
}
