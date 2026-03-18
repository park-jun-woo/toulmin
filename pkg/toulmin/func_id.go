//ff:func feature=engine type=util control=sequence
//ff:what funcID — extracts full path function name from function pointer via reflect
package toulmin

import (
	"fmt"
	"reflect"
	"runtime"
)

// funcID returns the full path name of a function from its pointer.
// e.g. "github.com/example/pkg.IsAdult"
// Returns a fallback string if runtime.FuncForPC returns nil.
func funcID(fn func(any, any) (bool, any)) string {
	ptr := reflect.ValueOf(fn).Pointer()
	f := runtime.FuncForPC(ptr)
	if f == nil {
		return fmt.Sprintf("unknown_%d", ptr)
	}
	return f.Name()
}
