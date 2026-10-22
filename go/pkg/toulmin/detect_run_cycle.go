//ff:func feature=engine type=validator control=sequence
//ff:what detectRunCycle — DFS over RunGraph edges to reject cyclic execution composition
package toulmin

import "fmt"

// detectRunCycle walks the graph-of-graphs reachable from root via each node's
// RunGraph edge, keyed by *Graph pointer identity, using a 3-color DFS
// (0=unvisited, 1=visiting, 2=done). Re-entering a visiting graph is an execution
// cycle and returns an error; a shared sub-graph reached by two paths (diamond DAG)
// is legal thanks to the done color. Execution composition must be a DAG.
func detectRunCycle(root *Graph) error {
	state := make(map[*Graph]int)
	var visit func(g *Graph) error
	visit = func(g *Graph) error {
		if state[g] == 1 {
			return fmt.Errorf("toulmin: run cycle detected at graph %q", g.name)
		}
		if state[g] == 2 {
			return nil
		}
		state[g] = 1
		for i := range g.rules {
			sub := g.rules[i].RunGraph
			if sub == nil {
				continue
			}
			if err := visit(sub); err != nil {
				return err
			}
		}
		state[g] = 2
		return nil
	}
	return visit(root)
}
