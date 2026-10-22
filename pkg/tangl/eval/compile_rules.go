//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what compileRules — compiles rules in dependency order into the compiled map
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// compileRules compiles each rule named in order (dependency-first) with the
// document resolver, writing the resulting closure into compiled so later rules
// can reference already-compiled siblings.
func compileRules(order []string, byName map[string]ast.InlineRule, dr docResolver, compiled map[string]evalFunc) error {
	for _, name := range order {
		ef, err := compileRule(byName[name], dr)
		if err != nil {
			return err
		}
		compiled[name] = ef
	}
	return nil
}
