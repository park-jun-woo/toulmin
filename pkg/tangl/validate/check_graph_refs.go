//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkGraphRefs — check that binding and eval graph references exist
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// checkGraphRefs checks that RuleBinding.Graph and EvalDecl.Graph reference existing GraphDecls.
func checkGraphRefs(f *parser.File) []string {
	var errs []string

	graphs := make(map[string]bool)
	for _, g := range f.Graphs {
		graphs[g.Name] = true
	}

	for _, b := range f.Bindings {
		if b.Graph != "" && !graphs[b.Graph] {
			errs = append(errs, fmt.Sprintf("binding %q references unknown graph %q (line %d)", b.Name, b.Graph, b.Line))
		}
	}

	for _, e := range f.Evals {
		if !graphs[e.Graph] {
			errs = append(errs, fmt.Sprintf("eval %q references unknown graph %q (line %d)", e.Name, e.Graph, e.Line))
		}
	}

	return errs
}
