//ff:func feature=engine type=util control=sequence
//ff:what funcID — extracts full path function name from function pointer via reflect
package toulmin

import (
	"reflect"
	"runtime"
)

// funcID returns the full path name of a function from its pointer.
// e.g. "github.com/example/pkg.IsAdult"
func funcID(fn func(any, any) (bool, any)) string {
	ptr := reflect.ValueOf(fn).Pointer()
	return runtime.FuncForPC(ptr).Name()
}
