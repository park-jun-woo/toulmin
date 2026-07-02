//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkExecFuncAliasRefs — verifies a do/undo exec's package alias is declared in See
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkExecFuncAliasRefs verifies that every do/undo exec's `<alias>.<name>`
// function reference alias was declared in the document's tangl:See section.
func checkExecFuncAliasRefs(doc *ast.Document) []error {
	aliases := seeAliasSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, e := range c.Execs {
			if e.Func == nil || e.Func.Alias == "" {
				continue
			}
			if !aliases[e.Func.Alias] {
				errs = append(errs, errAt(doc.Path, e.Line, "case %q: exec references undeclared package alias %q", c.Name, e.Func.Alias))
			}
		}
	}
	return errs
}
