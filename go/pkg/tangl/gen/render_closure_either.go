//ff:func feature=tangl type=codegen control=iteration dimension=1 pattern=recursive
//ff:what renderClosureEither — emits an Either group as a tanglOr closure literal
package gen

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureEither emits an Either group as a func() (bool, any) closure
// literal composing its term closures with tanglOr — Either semantics are
// identical to a Logic "or" (first-true short circuit), matching
// pkg/tangl/eval where both compile to the same or-closure. The parser
// produces exactly one term (the merged tree inside the "either" block);
// multiple terms join as or-alternatives defensively.
func renderClosureEither(e ast.Either, env map[string]binderInfo, depth int) (string, error) {
	terms := make([]string, 0, len(e.Terms))
	for _, t := range e.Terms {
		s, err := renderClosureExpr(t, env, depth)
		if err != nil {
			return "", err
		}
		terms = append(terms, s)
	}
	return "func() (bool, any) {\nreturn tanglOr(\n" + strings.Join(terms, ",\n") + ",\n)\n}", nil
}
