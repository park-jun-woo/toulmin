//ff:func feature=cli type=command control=iteration dimension=1
//ff:what formatEffects — renders effect closure entries as a column-aligned table
package tanglcli

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/park-jun-woo/toulmin/pkg/tangl/effects"
)

// formatEffects renders entries as a column-aligned table (kind, func,
// once, "(case / node)") matching the spec's effect summary example. An
// empty entries slice renders as an empty string.
func formatEffects(entries []effects.Entry) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', 0)
	for _, e := range entries {
		once := ""
		if e.Once {
			once = "once"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t(%s / %s)\n", e.Kind, refString(e.Func), once, e.Case, e.Node)
	}
	w.Flush()
	return buf.String()
}
