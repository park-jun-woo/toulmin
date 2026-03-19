//ff:func feature=engine type=util control=sequence
//ff:what ruleID — generates unique rule identifier from funcID and backing
package toulmin

import "fmt"

// ruleID returns a unique identifier for a rule.
// If backing is nil, returns funcID only.
// If backing is non-nil, returns funcID + "#" + backing string.
func ruleID(fn any, backing any) string {
	id := funcID(fn)
	if backing == nil {
		return id
	}
	return id + "#" + fmt.Sprint(backing)
}
