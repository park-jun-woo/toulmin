package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func renderStringListDef(w *strings.Builder, d ast.Definition, values []string) defInfo {
	ident := goIdent(d.Name) + "Spec"
	fmt.Fprintf(w, "var %s = toulmin.NewStringListSpec(%s, []string{", ident, strconv.Quote(d.Name))
	for i, value := range values {
		if i > 0 {
			w.WriteString(", ")
		}
		w.WriteString(strconv.Quote(value))
	}
	fmt.Fprintln(w, "})")
	fmt.Fprintln(w)
	return defInfo{Spec: ident, Kind: d.Kind}
}
