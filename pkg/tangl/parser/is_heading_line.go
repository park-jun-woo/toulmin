//ff:func feature=tangl type=parser control=sequence
//ff:what isHeadingLine — report whether a line is a markdown heading
package parser

import "strings"

// isHeadingLine reports whether line is a markdown heading of any level.
func isHeadingLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}
