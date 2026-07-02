//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkInternalCaseRefs — verifies Internal until/run/check targets name existing cases
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkInternalCaseRefs verifies that every case named by a tangl:Internal
// entry's `until`, `run`, or `check` clause exists in the document.
func checkInternalCaseRefs(doc *ast.Document) []error {
	cases := caseNameSet(doc)
	var errs []error
	for _, in := range doc.Internals {
		if in.Until != "" && !cases[in.Until] {
			errs = append(errs, errAt(doc.Path, in.Line, "tangl:Internal: until case %q does not exist", in.Until))
		}
		for _, name := range in.Runs {
			if !cases[name] {
				errs = append(errs, errAt(doc.Path, in.Line, "tangl:Internal: run case %q does not exist", name))
			}
		}
		for _, name := range in.Checks {
			if !cases[name] {
				errs = append(errs, errAt(doc.Path, in.Line, "tangl:Internal: check case %q does not exist", name))
			}
		}
	}
	return errs
}
