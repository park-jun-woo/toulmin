//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkFuncRefs — check that binding func references exist as imports or inline rules
package validate

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// checkFuncRefs checks that RuleBinding.Func references an existing Import alias or InlineRule name.
func checkFuncRefs(f *parser.File) []string {
	var errs []string

	aliases := make(map[string]bool)
	for _, imp := range f.Imports {
		aliases[imp.Alias] = true
	}

	inlineRules := make(map[string]bool)
	for _, r := range f.Rules {
		inlineRules[r.Name] = true
	}

	for _, b := range f.Bindings {
		errs = appendFuncRefError(errs, b, aliases, inlineRules)
	}

	return errs
}
