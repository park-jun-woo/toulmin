//ff:func feature=cli type=command control=iteration dimension=1
//ff:what runGraphGo — analyzes Go file for GraphBuilder defeat cycles
package cli

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/analyzer"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
	"github.com/spf13/cobra"
)

// runGraphGo analyzes a Go source file, extracts GraphBuilder defeat
// relationships via AST, and checks each graph for cycles.
func runGraphGo(cmd *cobra.Command, goPath string) error {
	graphs, err := analyzer.ExtractDefeats(goPath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", goPath, err)
	}
	if len(graphs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no GraphBuilder definitions found")
		return nil
	}
	hasError := false
	for _, g := range graphs {
		edges := make(map[string][]string)
		for _, d := range g.Defeats {
			edges[d.To] = append(edges[d.To], d.From)
		}
		if err := toulmin.DetectCycle(edges); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "graph %q: %v\n", g.Graph, err)
			hasError = true
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "graph %q: no cycles detected (%d rules, %d defeats)\n",
				g.Graph, len(g.Rules), len(g.Defeats))
		}
	}
	if hasError {
		return fmt.Errorf("cycle detected in one or more graphs")
	}
	return nil
}
