//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what equalStructIgnoringLine — field-wise struct equality that skips every field named Line
package tangeul

import "reflect"

// equalStructIgnoringLine compares two structs of identical type field by field,
// skipping any field named "Line" (the unstored source position) and recursing
// through equalIgnoringLine for the rest.
func equalStructIgnoringLine(x, y reflect.Value) bool {
	for i := 0; i < x.NumField(); i++ {
		if x.Type().Field(i).Name == "Line" {
			continue
		}
		if !equalIgnoringLine(x.Field(i), y.Field(i)) {
			return false
		}
	}
	return true
}
