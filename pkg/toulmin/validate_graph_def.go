//ff:func feature=graph type=validator control=iteration dimension=1
//ff:what ValidateGraphDef — validates GraphDef (graph name, rule names, roles, defeat edges, cycles)
package toulmin

import (
	"fmt"
	"strings"
)

// ValidateGraphDef checks that a GraphDef is well-formed:
// graph name is set, no duplicate rule names, roles are valid,
// all defeat edges reference existing rules, no cycles.
func ValidateGraphDef(def GraphDef) error {
	if def.Graph == "" {
		return fmt.Errorf("graph name is required")
	}
	known := make(map[string]bool)
	var errs []string
	for _, r := range def.Rules {
		if known[r.Name] {
			errs = append(errs, fmt.Sprintf("  duplicate rule name %q", r.Name))
		}
		known[r.Name] = true
		switch r.Role {
		case "rule", "counter", "except":
		default:
			errs = append(errs, fmt.Sprintf("  rule %q has unknown role %q", r.Name, r.Role))
		}
	}
	for _, d := range def.Defeats {
		if !known[d.From] {
			errs = append(errs, fmt.Sprintf("  defeats from unknown rule %q", d.From))
		}
		if !known[d.To] {
			errs = append(errs, fmt.Sprintf("  defeats to unknown rule %q", d.To))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("graph %q validation failed:\n%s", def.Graph, strings.Join(errs, "\n"))
	}
	edges := make(map[string][]string)
	for _, d := range def.Defeats {
		edges[d.To] = append(edges[d.To], d.From)
	}
	if err := DetectCycle(edges); err != nil {
		return fmt.Errorf("graph %q: %w", def.Graph, err)
	}
	return nil
}
