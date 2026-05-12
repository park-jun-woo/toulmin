//ff:func feature=engine type=util control=sequence
//ff:what shortName — extracts short name from full path string
package toulmin

import "strings"

// shortName returns the short name from a full path.
// e.g. "github.com/example/pkg.IsAdult" → "IsAdult"
func shortName(full string) string {
	base, spec, hasSpec := strings.Cut(full, "#")
	if idx := strings.LastIndex(base, "."); idx >= 0 {
		base = base[idx+1:]
	}
	if hasSpec {
		return base + "#" + spec
	}
	return base
}
