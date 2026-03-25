//ff:func feature=tangl type=parser control=sequence
//ff:what isBinding — check if a line contains a rule binding pattern
package parser

import "strings"

// isBinding checks if a line matches a rule binding pattern.
func isBinding(line string) bool {
	return strings.Contains(line, " is a rule ") ||
		strings.Contains(line, " is a counter ") ||
		strings.Contains(line, " is an except ")
}
