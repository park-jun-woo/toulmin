//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkEndpointCaseRefs — verifies provides run/check targets name existing cases
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkEndpointCaseRefs verifies that every `run`/`check` case named by a
// tangl:Provides endpoint exists in the document.
func checkEndpointCaseRefs(doc *ast.Document) []error {
	cases := caseNameSet(doc)
	var errs []error
	for _, ep := range doc.Provides {
		for _, name := range ep.Runs {
			if !cases[name] {
				errs = append(errs, errAt(doc.Path, ep.Line, "provides %q: run case %q does not exist", ep.Name, name))
			}
		}
		for _, name := range ep.Checks {
			if !cases[name] {
				errs = append(errs, errAt(doc.Path, ep.Line, "provides %q: check case %q does not exist", ep.Name, name))
			}
		}
	}
	return errs
}
