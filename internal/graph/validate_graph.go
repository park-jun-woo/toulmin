//ff:func feature=graph type=validator control=iteration dimension=1
//ff:what ValidateGraph — validates all defeats references in rule set
package graph

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// ValidateGraph checks that all defeats targets reference existing rules.
// Returns nil if valid, or an error listing all invalid references.
func ValidateGraph(metas []toulmin.RuleMeta) error {
	known := make(map[string]bool)
	for _, m := range metas {
		known[m.Name] = true
	}
	var errs []string
	for _, m := range metas {
		errs = append(errs, validateDefeats(m, known)...)
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("graph validation failed:\n%s", strings.Join(errs, "\n"))
}
