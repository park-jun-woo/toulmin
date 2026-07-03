//ff:func feature=tangl type=parser control=sequence
//ff:what isHeading — report whether a line's text is a markdown heading of any level
package scan

import "strings"

// isHeading reports whether text is a markdown heading of any level, matching
// the parser's isHeadingLine rule: a heading closes a section's scope.
func isHeading(text string) bool {
	return strings.HasPrefix(strings.TrimSpace(text), "#")
}
