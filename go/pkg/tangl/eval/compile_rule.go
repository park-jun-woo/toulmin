//ff:func feature=tangl type=engine control=sequence
//ff:what compileRule — compile an expression rule or bind a host predicate declaration
package eval

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func compileRule(rule ast.InlineRule, res Resolver) (evalFunc, error) {
	if rule.Using == nil {
		return compileExpr(rule.Cond, res)
	}
	if res == nil {
		return nil, fmt.Errorf("rule %q: no resolver for host predicate %s.%s", rule.Name, rule.Using.Alias, rule.Using.Name)
	}
	fn, ok := res.Predicate(rule.Using.Alias, rule.Using.Name)
	if !ok {
		return nil, fmt.Errorf("rule %q: unresolved host predicate %s.%s", rule.Name, rule.Using.Alias, rule.Using.Name)
	}
	return func(sc *scope) (bool, any) { return fn(sc.ctx, sc.specs) }, nil
}
