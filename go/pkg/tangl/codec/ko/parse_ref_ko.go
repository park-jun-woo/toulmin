//ff:func feature=tangl type=parser control=sequence
//ff:what parseRefKo — parse a FUNC_REF (`alias`.`name` or local `name`)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseRefKo parses a FUNC_REF: either a qualified `alias`.`name` reference or a
// bare local `name`. It mirrors pkg/tangl/parser.parseRef; the KO package
// redefines it so definition/type tails resolve backtick references locally.
func parseRefKo(s string) (ast.Ref, string, bool) {
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
