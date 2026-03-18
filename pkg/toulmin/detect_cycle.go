//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what detectCycle — DFS cycle detection on directed defeat graph
package toulmin

import "fmt"

// detectCycle checks for cycles in a directed graph using DFS.
// Returns an error naming the node where a cycle was detected.
func detectCycle(edges map[string][]string) error {
	state := make(map[string]int) // 0=unvisited, 1=visiting, 2=done
	var visit func(id string) error
	visit = func(id string) error {
		if state[id] == 1 {
			return fmt.Errorf("cycle detected at %q", id)
		}
		if state[id] == 2 {
			return nil
		}
		state[id] = 1
		for _, aid := range edges[id] {
			if err := visit(aid); err != nil {
				return err
			}
		}
		state[id] = 2
		return nil
	}
	for id := range edges {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
