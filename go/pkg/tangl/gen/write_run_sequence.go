//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeRunSequence — writes one Run+compensate block per endpoint run case
package gen

import (
	"fmt"
	"strings"
)

// writeRunSequence writes one Run+compensate block per case in runs, in
// document order, sharing the single InitCompensation armed by the
// caller: the first handler error triggers LIFO Compensate (escalating to
// Review on a failed compensation) and returns immediately.
func writeRunSequence(w *strings.Builder, runs []string, returnsResults bool) {
	zero := "err"
	if returnsResults {
		zero = "nil, err"
	}
	for _, c := range runs {
		fmt.Fprintf(w, "\tif _, _, err := %sGraph.Run(ctx); err != nil {\n", goIdent(c))
		fmt.Fprintln(w, "\t\tif cerr := tangl.Compensate(ctx); cerr != nil {")
		if returnsResults {
			fmt.Fprintln(w, "\t\t\treturn nil, tangl.Review(ctx, err, cerr)")
		} else {
			fmt.Fprintln(w, "\t\t\treturn tangl.Review(ctx, err, cerr)")
		}
		fmt.Fprintln(w, "\t\t}")
		fmt.Fprintf(w, "\t\treturn %s\n", zero)
		fmt.Fprintln(w, "\t}")
	}
}
