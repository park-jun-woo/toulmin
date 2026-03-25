//ff:func feature=tangl type=validator control=sequence
//ff:what appendFuncRefError — classify a binding func reference and append error if unresolved
package validate

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// appendFuncRefError checks a single binding's func reference and appends an error if unresolved.
func appendFuncRefError(errs []string, b parser.RuleBinding, aliases map[string]bool, inlineRules map[string]bool) []string {
	funcRef := b.Func
	if strings.Contains(funcRef, ".") {
		parts := strings.SplitN(funcRef, ".", 2)
		if !aliases[parts[0]] {
			errs = append(errs, fmt.Sprintf("binding %q references unknown import alias %q (line %d)", b.Name, parts[0], b.Line))
		}
		return errs
	}
	if !inlineRules[funcRef] && !aliases[funcRef] {
		errs = append(errs, fmt.Sprintf("binding %q references unknown function %q (line %d)", b.Name, funcRef, b.Line))
	}
	return errs
}
