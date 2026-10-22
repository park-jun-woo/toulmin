//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what structField — finds an exported struct field by exact then case-insensitive name
package eval

import (
	"reflect"
	"strings"
)

// structField returns the value of the exported field of v matching name,
// trying an exact match first and then a case-insensitive match (so a backtick
// step like `trustScore` resolves the exported field TrustScore). It reports
// ok=false when no exported field matches.
func structField(v reflect.Value, name string) (reflect.Value, bool) {
	t := v.Type()
	if f, ok := t.FieldByName(name); ok && f.PkgPath == "" {
		return v.FieldByIndex(f.Index), true
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath == "" && strings.EqualFold(sf.Name, name) {
			return v.Field(i), true
		}
	}
	return reflect.Value{}, false
}
