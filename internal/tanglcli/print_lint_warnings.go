//ff:func feature=cli type=command control=iteration dimension=1
//ff:what printLintWarnings — writes each lint warning to w, one per line
package tanglcli

import (
	"fmt"
	"io"
)

// printLintWarnings writes each warning to w, one per line.
func printLintWarnings(w io.Writer, warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintln(w, warning)
	}
}
