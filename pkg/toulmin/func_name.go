//ff:func feature=engine type=util control=sequence
//ff:what FuncName — extracts function name from function pointer via reflect
package toulmin

import (
	"reflect"
	"runtime"
	"strings"
)

// FuncName returns the short name of a function from its pointer.
// e.g. "github.com/example/pkg.IsAdult" → "IsAdult"
func FuncName(fn func(any, any) bool) string {
	ptr := reflect.ValueOf(fn).Pointer()
	full := runtime.FuncForPC(ptr).Name()
	if idx := strings.LastIndex(full, "."); idx >= 0 {
		return full[idx+1:]
	}
	return full
}
