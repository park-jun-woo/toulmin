//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeImportBlock — writes a parenthesized Go import block
package gen

import (
	"fmt"
	"strconv"
	"strings"
)

// writeImportBlock writes a parenthesized Go import block from specs, one
// import per line, emitting an explicit alias only when one was recorded
// for that importSpec.
func writeImportBlock(w *strings.Builder, specs []importSpec) {
	fmt.Fprintln(w, "import (")
	for _, s := range specs {
		if s.Alias != "" {
			fmt.Fprintf(w, "\t%s %s\n", s.Alias, strconv.Quote(s.Path))
		} else {
			fmt.Fprintf(w, "\t%s\n", strconv.Quote(s.Path))
		}
	}
	fmt.Fprintln(w, ")")
}
