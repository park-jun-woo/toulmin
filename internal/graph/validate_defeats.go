//ff:func feature=graph type=validator control=iteration dimension=1
//ff:what validateDefeats — checks defeats targets exist for a single rule
package graph

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// validateDefeats returns error messages for defeats targets not in known set.
func validateDefeats(meta toulmin.RuleMeta, known map[string]bool) []string {
	var errs []string
	for _, target := range meta.Defeats {
		if !known[target] {
			errs = append(errs, fmt.Sprintf("  rule %q defeats unknown rule %q", meta.Name, target))
		}
	}
	return errs
}
