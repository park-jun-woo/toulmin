//ff:func feature=graph type=validator control=iteration dimension=1
//ff:what Validate — validates GraphDef (rule names exist, defeats targets valid, no cycles)
package graphdef

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Validate checks that a GraphDef is well-formed:
// graph name is set, all defeat edges reference existing rules, no cycles.
func Validate(def *GraphDef) error {
	if def.Graph == "" {
		return fmt.Errorf("graph name is required")
	}
	known := make(map[string]bool)
	for _, r := range def.Rules {
		known[r.Name] = true
	}
	var errs []string
	for _, d := range def.Defeats {
		if !known[d.From] {
			errs = append(errs, fmt.Sprintf("  defeats from unknown rule %q", d.From))
		}
		if !known[d.To] {
			errs = append(errs, fmt.Sprintf("  defeats to unknown rule %q", d.To))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("graph validation failed:\n%s", strings.Join(errs, "\n"))
	}
	edges := buildEdgesFromDef(def.Defeats)
	if err := toulmin.DetectCycle(edges); err != nil {
		return fmt.Errorf("graph %q: %w", def.Graph, err)
	}
	return nil
}
