//ff:func feature=engine type=util control=iteration dimension=1
//ff:what ruleID — generates unique rule identifier from funcID and specs
package toulmin

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ruleID returns a unique identifier for a rule.
// If specs is empty, returns funcID only.
// If specs is non-empty, returns funcID + "#" + sorted spec identifiers,
// where each identifier combines the spec's type name (SpecName) with a
// JSON value serialization of its fields. JSON serializes nested pointers
// and interface values by value, so equal specs produce equal identifiers
// regardless of instance address. Specs that fail to serialize fall back to
// %+v (collision-avoidance preferred over stability).
func ruleID(fn any, specs Specs) string {
	id := funcID(fn)
	if len(specs) == 0 {
		return id
	}
	names := make([]string, len(specs))
	for i, s := range specs {
		b, err := json.Marshal(s)
		if err != nil {
			b = []byte(fmt.Sprintf("%+v", s))
		}
		names[i] = s.SpecName() + string(b)
	}
	sort.Strings(names)
	return id + "#" + strings.Join(names, "+")
}
