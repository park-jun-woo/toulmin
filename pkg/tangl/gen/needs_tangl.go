//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what needsTangl — reports whether the generated file calls the tangl runtime
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// needsTangl reports whether any generated code will call the tangl
// runtime: every Provides "run" endpoint and every Internal trigger
// always wraps its pass in the Init/Compensate/Commit cycle, and any
// "once" or "undo" edge needs the once-guard or compensation-stack
// helpers too.
func needsTangl(doc *ast.Document) bool {
	for _, ep := range doc.Provides {
		if len(ep.Requires) > 0 || len(ep.Runs) > 0 {
			return true
		}
	}
	if len(doc.Internals) > 0 {
		return true
	}
	for _, c := range doc.Cases {
		for _, e := range c.Execs {
			if e.Once || e.Kind == ast.UndoExec {
				return true
			}
		}
	}
	return false
}
