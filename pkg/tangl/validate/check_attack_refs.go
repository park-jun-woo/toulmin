//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkAttackRefs — check that attack attacker and target reference existing bindings
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// checkAttackRefs checks that AttackDecl.Attacker and .Target reference existing RuleBinding names.
func checkAttackRefs(f *parser.File) []string {
	var errs []string

	bindings := make(map[string]bool)
	for _, b := range f.Bindings {
		bindings[b.Name] = true
	}

	for _, a := range f.Attacks {
		if !bindings[a.Attacker] {
			errs = append(errs, fmt.Sprintf("attack attacker %q references unknown binding (line %d)", a.Attacker, a.Line))
		}
		if !bindings[a.Target] {
			errs = append(errs, fmt.Sprintf("attack target %q references unknown binding (line %d)", a.Target, a.Line))
		}
	}

	return errs
}
