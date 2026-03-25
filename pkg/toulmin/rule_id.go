//ff:func feature=engine type=util control=iteration dimension=1
//ff:what ruleID — generates unique rule identifier from funcID and specs
package toulmin

import (
	"fmt"
	"sort"
	"strings"
)

// ruleID returns a unique identifier for a rule.
// If specs is empty, returns funcID only.
// If specs is non-empty, returns funcID + "#" + sorted spec names.
func ruleID(fn any, specs Specs) string {
	id := funcID(fn)
	if len(specs) == 0 {
		return id
	}
	names := make([]string, len(specs))
	for i, s := range specs {
		names[i] = fmt.Sprint(s)
	}
	sort.Strings(names)
	return id + "#" + strings.Join(names, "+")
}
