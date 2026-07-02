//ff:func feature=tangl type=codegen control=sequence
//ff:what renderCompareExpr — compiles a Compare leaf into a tanglCompare call
package gen

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderCompareExpr compiles a Compare leaf into a tanglCompare(...)
// call, quoting its field, operator, and value as Go string literals.
func renderCompareExpr(c ast.Compare) string {
	return fmt.Sprintf("tanglCompare(ctx, %s, %s, %s)", strconv.Quote(c.Field), strconv.Quote(c.Op), strconv.Quote(c.Value))
}
