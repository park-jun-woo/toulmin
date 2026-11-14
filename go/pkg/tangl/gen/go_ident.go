//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what goIdent — converts a tangl backtick name to an unexported camelCase Go identifier
package gen

import "strings"

// goIdent converts a tangl backtick name to an unexported Go identifier:
// space-separated words become camelCase ("order received" becomes
// "orderReceived"). A single-word name passes through sanitized but
// otherwise unchanged.
func goIdent(name string) string {
	words := strings.Fields(name)
	var b strings.Builder
	for i, w := range words {
		w = sanitizeWord(w)
		if w == "" {
			continue
		}
		if i == 0 {
			b.WriteString(strings.ToLower(w[:1]) + w[1:])
		} else {
			b.WriteString(strings.ToUpper(w[:1]) + w[1:])
		}
	}
	out := b.String()
	if out == "" {
		return "_"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "_" + out
	}
	return out
}
