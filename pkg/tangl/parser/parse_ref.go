//ff:func feature=tangl type=parser control=sequence
//ff:what parseRef — parse a FUNC_REF (`alias`.`name` or local `name`)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseRef parses a FUNC_REF: either a qualified `alias`.`name` reference or
// a bare local `name` reference.
func parseRef(s string) (ast.Ref, string, bool) {
	first, rest, ok := takeBacktick(s)
	if !ok {
		return ast.Ref{}, s, false
	}
	if strings.HasPrefix(rest, ".") {
		second, rest2, ok2 := takeBacktick(rest[1:])
		if !ok2 {
			return ast.Ref{}, s, false
		}
		return ast.Ref{Alias: first, Name: second}, rest2, true
	}
	return ast.Ref{Name: first}, rest, true
}
