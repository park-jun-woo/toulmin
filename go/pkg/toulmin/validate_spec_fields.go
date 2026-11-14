//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what validateSpecFields — validates that Spec struct has no func fields
package toulmin

import (
	"fmt"
	"reflect"
)

// validateSpecFields checks all fields of the spec struct via reflect.
// Returns an error if any field has Kind == reflect.Func.
func validateSpecFields(s Spec) error {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Func {
			return fmt.Errorf("spec %s has func field %q — func fields are not allowed",
				s.SpecName(), t.Field(i).Name)
		}
	}
	return nil
}
