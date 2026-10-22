//ff:func feature=tangl type=util control=sequence
//ff:what normalizeLoc — replace each "path:line:" message prefix with a fixed LOC: token
package tangeul

import "regexp"

var locPrefix = regexp.MustCompile(`(?m)^.*?:\d+:`)

// normalizeLoc rewrites the leading "path:line:" of every line to the fixed token
// "LOC:", so validate/Lint messages from the two routes compare on their semantic
// content only. The text route emits the .md path with real line numbers; the
// .geul route emits the .geul path with Line 0 — both non-content metadata
// the stream never carried (README contract §2, Phase 006).
func normalizeLoc(s string) string {
	return locPrefix.ReplaceAllString(s, "LOC:")
}
