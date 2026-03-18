//ff:func feature=engine type=util control=sequence
//ff:what shortName — extracts short name from full path string
package toulmin

import "strings"

// shortName returns the short name from a full path.
// e.g. "github.com/example/pkg.IsAdult" → "IsAdult"
func shortName(full string) string {
	if idx := strings.LastIndex(full, "."); idx >= 0 {
		return full[idx+1:]
	}
	return full
}
