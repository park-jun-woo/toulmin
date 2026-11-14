//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeTickRunSequence — writes one Run+compensate block per ticker run case
package gen

import (
	"fmt"
	"strings"
)

// writeTickRunSequence writes one Run+compensate block per case in runs,
// indented for the ticker's per-tick body: unlike an endpoint wrapper, a
// failed compensation logs via tangl.Review and the loop simply
// continues to the next tick rather than returning — Internal has no
// caller to report an error to.
func writeTickRunSequence(w *strings.Builder, runs []string) {
	for _, c := range runs {
		fmt.Fprintf(w, "\t\tif _, _, err := %sGraph.Run(ctx); err != nil {\n", goIdent(c))
		fmt.Fprintln(w, "\t\t\tif cerr := tangl.Compensate(ctx); cerr != nil {")
		fmt.Fprintln(w, "\t\t\t\ttangl.Review(ctx, err, cerr)")
		fmt.Fprintln(w, "\t\t\t}")
		fmt.Fprintln(w, "\t\t\tcontinue")
		fmt.Fprintln(w, "\t\t}")
	}
}
