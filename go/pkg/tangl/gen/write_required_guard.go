//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeRequiredGuard — writes the tangl.Required precondition check
package gen

import (
	"fmt"
	"strconv"
	"strings"
)

// writeRequiredGuard writes the tangl.Required(...) precondition check,
// if fields is non-empty, returning the mode-appropriate zero value on
// error.
func writeRequiredGuard(w *strings.Builder, fields []string, returnsResults bool) {
	if len(fields) == 0 {
		return
	}
	args := make([]string, len(fields))
	for i, f := range fields {
		args[i] = strconv.Quote(f)
	}
	fmt.Fprintf(w, "\tif err := tangl.Required(ctx, %s); err != nil {\n", strings.Join(args, ", "))
	if returnsResults {
		fmt.Fprintln(w, "\t\treturn nil, err")
	} else {
		fmt.Fprintln(w, "\t\treturn err")
	}
	fmt.Fprintln(w, "\t}")
}
