//ff:func feature=tangl type=codegen control=sequence
//ff:what tickFuncName — builds a unique unexported name for an every-ticker runner
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// tickFuncName builds a unique unexported name for an "every" ticker
// runner, seeded from its first Run/Check case name (falling back to
// "tick") with the Internal's document index appended for uniqueness.
func tickFuncName(in ast.Internal, idx int) string {
	seed := "tick"
	if len(in.Runs) > 0 {
		seed = in.Runs[0]
	} else if len(in.Checks) > 0 {
		seed = in.Checks[0]
	}
	return fmt.Sprintf("tick%s%d", goIdentExported(seed), idx)
}
