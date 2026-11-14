//ff:func feature=cli type=command control=sequence
//ff:what trimBackticks — strips leading and trailing backtick characters from a name
package tanglcli

import "strings"

// trimBackticks strips any leading and trailing backtick characters from
// s, so an endpoint name may be given with or without its surrounding
// backticks.
func trimBackticks(s string) string {
	return strings.Trim(s, "`")
}
