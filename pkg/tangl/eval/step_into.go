//ff:func feature=tangl type=engine control=sequence
//ff:what stepInto — resolves one path step into cur via map key or struct field
package eval

import "reflect"

// stepInto resolves one step into cur: a map[string]any key lookup, else a
// struct field access (auto-dereferencing a pointer first). A nil value,
// missing key, non-struct, or missing field yields ok=false.
func stepInto(cur any, step string) (any, bool) {
	if cur == nil {
		return nil, false
	}
	if m, ok := cur.(map[string]any); ok {
		v, present := m[step]
		return v, present
	}
	rv := reflect.ValueOf(cur)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, false
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, false
	}
	fv, ok := structField(rv, step)
	if !ok {
		return nil, false
	}
	return fv.Interface(), true
}
