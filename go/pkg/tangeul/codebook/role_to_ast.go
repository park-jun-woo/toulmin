//ff:func feature=tangl type=util control=selection
//ff:what RoleToAST — map a RoleCode back to its ast.Role, erroring on unknown
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// RoleToAST maps a codebook RoleCode back to its ast.Role (§2.2), the inverse
// of RoleFromAST. An unrecognized code is a hard error — never a silent
// fallback.
func RoleToAST(c toulmin.RoleCode) (ast.Role, error) {
	switch c {
	case toulmin.RoleGeneral:
		return ast.GeneralRule, nil
	case toulmin.RoleCounter:
		return ast.CounterRule, nil
	case toulmin.RoleExcept:
		return ast.ExceptRule, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown toulmin.RoleCode %d", uint16(c))
	}
}
