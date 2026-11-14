//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what renderClosureLogic — emits a Logic node as a tanglAnd/tanglOr closure literal
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureLogic emits a Logic node as a func() (bool, any) closure
// literal composing its term closures with the short-circuit tanglAnd
// ("and", first false wins) or tanglOr ("or", first true wins) helper,
// recursing into each term via renderClosureExpr.
func renderClosureLogic(l ast.Logic, env map[string]binderInfo, depth int) (string, error) {
	helper := "tanglAnd"
	if l.Op == "or" {
		helper = "tanglOr"
	}
	terms := make([]string, 0, len(l.Terms))
	for _, t := range l.Terms {
		s, err := renderClosureExpr(t, env, depth)
		if err != nil {
			return "", err
		}
		terms = append(terms, s)
	}
	return "func() (bool, any) {\nreturn " + helper + "(\n" + strings.Join(terms, ",\n") + ",\n)\n}", nil
}
