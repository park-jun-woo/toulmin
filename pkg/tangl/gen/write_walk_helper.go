//ff:func feature=tangl type=codegen control=sequence
//ff:what writeWalkHelper — writes the shared tanglWalk and tanglField functions
package gen

import "strings"

// writeWalkHelper writes the shared tanglWalk function and its tanglField
// companion: tanglWalk walks map-key/struct-field steps into a value
// (auto-dereferencing pointers), and tanglField finds an exported struct
// field by exact then case-insensitive name — mirroring
// pkg/tangl/eval.walkSteps/stepInto/structField. A missing key/field or nil
// intermediate yields (nil, false).
func writeWalkHelper(w *strings.Builder) {
	w.WriteString(`func tanglWalk(cur any, steps ...string) (any, bool) {
	for _, step := range steps {
		if cur == nil {
			return nil, false
		}
		if m, isMap := cur.(map[string]any); isMap {
			v, present := m[step]
			if !present {
				return nil, false
			}
			cur = v
			continue
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
		fv, ok := tanglField(rv, step)
		if !ok {
			return nil, false
		}
		cur = fv.Interface()
	}
	return cur, true
}

func tanglField(v reflect.Value, name string) (reflect.Value, bool) {
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

`)
}
