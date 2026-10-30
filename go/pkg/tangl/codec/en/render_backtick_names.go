//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what renderBacktickNames — render names as backtick-quoted words joined by " and "
package en

import "strings"

// renderBacktickNames renders each name as "`name`" and joins them with " and "
// (the grounded-on list of a host rule declaration).
func renderBacktickNames(names []string) string {
	parts := make([]string, len(names))
	for i, name := range names {
		parts[i] = "`" + name + "`"
	}
	return strings.Join(parts, " and ")
}
