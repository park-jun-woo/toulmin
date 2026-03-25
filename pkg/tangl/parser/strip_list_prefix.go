//ff:func feature=tangl type=parser control=sequence
//ff:what stripListPrefix — remove markdown list prefix from a line
package parser

import "strings"

// stripListPrefix removes "- " or "N. " prefix from a markdown list item.
func stripListPrefix(s string) string {
	if strings.HasPrefix(s, "- ") {
		return s[2:]
	}
	return stripNumberedPrefix(s)
}
