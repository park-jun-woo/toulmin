//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what quoteSteps — renders path steps as a comma-separated Go string-literal list
package gen

import (
	"strconv"
	"strings"
)

// quoteSteps renders path steps as a comma-separated list of Go string
// literals, the argument form the tanglCtxVal/tanglSpecVal/tanglWalk helpers
// take.
func quoteSteps(steps []string) string {
	parts := make([]string, 0, len(steps))
	for _, s := range steps {
		parts = append(parts, strconv.Quote(s))
	}
	return strings.Join(parts, ", ")
}
