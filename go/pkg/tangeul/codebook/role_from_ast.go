//ff:func feature=tangl type=util control=selection
//ff:what RoleFromAST — map an ast.Role to its 2-bit RoleCode, erroring on unknown
package codebook

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// RoleFromAST maps an ast.Role to its codebook RoleCode (§2.2). An unrecognized
// Role is a hard error — never a silent fallback.
func RoleFromAST(r ast.Role) (toulmin.RoleCode, error) {
	switch r {
	case ast.GeneralRule:
		return toulmin.RoleGeneral, nil
	case ast.CounterRule:
		return toulmin.RoleCounter, nil
	case ast.ExceptRule:
		return toulmin.RoleExcept, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown ast.Role %d", int(r))
	}
}
