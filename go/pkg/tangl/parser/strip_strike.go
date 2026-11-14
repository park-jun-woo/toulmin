//ff:func feature=tangl type=parser control=sequence
//ff:what stripStrike — detect and unwrap a `~~...~~` strikethrough comment item
package parser

import "strings"

// stripStrike reports whether content is wrapped in `~~...~~` (a
// strikethrough comment item, per spec) and returns the unwrapped text.
func stripStrike(content string) (text string, struck bool) {
	if strings.HasPrefix(content, "~~") && strings.HasSuffix(content, "~~") && len(content) >= 4 {
		return content[2 : len(content)-2], true
	}
	return content, false
}
