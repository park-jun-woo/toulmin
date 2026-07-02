//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseItemLine — split a trimmed line into its list marker and content
package parser

import "strconv"

// parseItemLine recognizes a `- ` (unordered) or `N. ` (ordered) list marker
// at the start of trimmed and returns the remaining content. ok is false for
// any line that is not a list item (ignored by the caller).
func parseItemLine(trimmed string) (content string, ordered bool, number int, ok bool) {
	if len(trimmed) >= 2 && trimmed[0] == '-' && trimmed[1] == ' ' {
		return trimStart(trimmed[2:]), false, 0, true
	}
	i := 0
	for i < len(trimmed) && trimmed[i] >= '0' && trimmed[i] <= '9' {
		i++
	}
	if i == 0 || i+1 >= len(trimmed) || trimmed[i] != '.' || trimmed[i+1] != ' ' {
		return "", false, 0, false
	}
	n, err := strconv.Atoi(trimmed[:i])
	if err != nil {
		return "", false, 0, false
	}
	return trimStart(trimmed[i+2:]), true, n, true
}
