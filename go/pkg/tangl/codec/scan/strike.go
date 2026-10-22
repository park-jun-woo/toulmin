//ff:func feature=tangl type=parser control=sequence
//ff:what strike — detect and unwrap a `~~...~~` strikethrough comment item
package scan

import "strings"

// strike reports whether content is wrapped in `~~...~~` (a strikethrough
// comment item, per spec) and returns the unwrapped text, matching the parser's
// stripStrike. The consumer lowers struck items to PROSE per the v0.1 rule.
func strike(content string) (text string, struck bool) {
	if strings.HasPrefix(content, "~~") && strings.HasSuffix(content, "~~") && len(content) >= 4 {
		return content[2 : len(content)-2], true
	}
	return content, false
}
