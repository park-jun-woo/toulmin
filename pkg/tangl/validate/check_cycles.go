//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what checkCycles — detect cycles in attack edges using DFS
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// checkCycles builds attack edges and uses DFS to detect cycles.
func checkCycles(f *parser.File) []string {
	adj := make(map[string][]string)
	for _, a := range f.Attacks {
		adj[a.Attacker] = append(adj[a.Attacker], a.Target)
	}

	if len(adj) == 0 {
		return nil
	}

	white := 0
	gray := 1
	black := 2
	color := make(map[string]int)

	var errs []string
	var dfs func(node string) bool
	dfs = func(node string) bool {
		color[node] = gray
		for _, next := range adj[node] {
			if color[next] == gray {
				errs = append(errs, fmt.Sprintf("cycle detected: %s -> %s", node, next))
				return true
			}
			if color[next] == white {
				if dfs(next) {
					return true
				}
			}
		}
		color[node] = black
		return false
	}

	for node := range adj {
		if color[node] == white {
			dfs(node)
		}
	}

	_ = white
	return errs
}
