//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkRunCaseRefs — verifies run exec targets name an existing case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkRunCaseRefs verifies that every `run <case> when <node>` edge's
// target case exists in the document.
func checkRunCaseRefs(doc *ast.Document) []error {
	cases := caseNameSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, e := range c.Execs {
			if e.Kind != ast.RunExec {
				continue
			}
			if !cases[e.Case] {
				errs = append(errs, errAt(doc.Path, e.Line, "case %q: run target case %q does not exist", c.Name, e.Case))
			}
		}
	}
	return errs
}
