//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what validateBackingFields — validates that Backing struct has no func fields
package toulmin

import (
	"fmt"
	"reflect"
)

// validateBackingFields checks all fields of the backing struct via reflect.
// Returns an error if any field has Kind == reflect.Func.
func validateBackingFields(b Backing) error {
	t := reflect.TypeOf(b)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Func {
			return fmt.Errorf("backing %s has func field %q — func fields are not allowed",
				b.BackingName(), t.Field(i).Name)
		}
	}
	return nil
}
