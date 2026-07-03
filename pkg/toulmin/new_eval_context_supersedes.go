//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what detectSupersedeCycle — rejects cyclic supersede specs independently of the Attacks graph
package toulmin

import "fmt"

// detectSupersedeCycle builds a from→to adjacency map from the supersede edges and
// runs DetectCycle on it, independently of the Attacks (defeats) graph. Supersede
// edges never enter the h-Categoriser traversal, so a supersede cycle cannot cause
// non-termination; this gate instead rejects mutually-excluding specs (X excludes Y /
// Y excludes X) as an authoring error, keeping the deterministic priority well-formed.
// The returned error is wrapped so its message reads "toulmin: supersede cycle detected at %q".
func detectSupersedeCycle(supersedes []defeatEdge) error {
	adj := make(map[string][]string)
	for _, s := range supersedes {
		adj[s.from] = append(adj[s.from], s.to)
	}
	if err := DetectCycle(adj); err != nil {
		return fmt.Errorf("toulmin: supersede %w", err)
	}
	return nil
}
