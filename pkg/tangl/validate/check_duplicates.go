//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkDuplicates — check for duplicate graph, inline rule, and binding names
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// checkDuplicates checks for duplicate names across graphs, inline rules, and bindings.
func checkDuplicates(f *parser.File) []string {
	var errs []string

	graphNames := make(map[string]int)
	for _, g := range f.Graphs {
		if prev, ok := graphNames[g.Name]; ok {
			errs = append(errs, fmt.Sprintf("duplicate graph name %q (lines %d and %d)", g.Name, prev, g.Line))
		}
		graphNames[g.Name] = g.Line
	}

	ruleNames := make(map[string]int)
	for _, r := range f.Rules {
		if prev, ok := ruleNames[r.Name]; ok {
			errs = append(errs, fmt.Sprintf("duplicate inline rule name %q (lines %d and %d)", r.Name, prev, r.Line))
		}
		ruleNames[r.Name] = r.Line
	}

	bindingNames := make(map[string]int)
	for _, b := range f.Bindings {
		if prev, ok := bindingNames[b.Name]; ok {
			errs = append(errs, fmt.Sprintf("duplicate binding name %q (lines %d and %d)", b.Name, prev, b.Line))
		}
		bindingNames[b.Name] = b.Line
	}

	return errs
}
