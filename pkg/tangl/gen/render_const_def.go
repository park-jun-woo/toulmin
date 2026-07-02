//ff:func feature=tangl type=codegen control=sequence
//ff:what renderConstDef — writes one ConstDef as a Go const plus an optional Spec var
package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderConstDef writes one ConstDef as a Go const (a numeric literal is
// emitted as-is; a non-numeric literal like a unit-suffixed value is
// quoted as a Go string) and, if the definition carries a SpecRef, an
// accompanying package-level Spec var built by calling the referenced
// constructor with that constant.
func renderConstDef(w *strings.Builder, d ast.Definition) defInfo {
	ident := goIdent(d.Name)
	if isNumericLiteral(d.Value) {
		fmt.Fprintf(w, "const %s = %s\n", ident, d.Value)
	} else {
		fmt.Fprintf(w, "const %s = %s\n", ident, strconv.Quote(d.Value))
	}
	info := defInfo{Const: ident, Kind: ast.ConstDef}
	if d.SpecRef != nil {
		specIdent := ident + "Spec"
		fmt.Fprintf(w, "var %s = %s(%s)\n", specIdent, refSelector(d.SpecRef), ident)
		info.Spec = specIdent
	}
	fmt.Fprintln(w)
	return info
}
