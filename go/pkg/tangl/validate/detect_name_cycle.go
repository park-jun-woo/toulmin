//ff:func feature=tangl type=validator control=iteration dimension=1
//ff:what detectNameCycle — DFS cycle detection over a named-edge graph
package validate

// detectNameCycle performs a 3-color DFS (0=unvisited, 1=visiting, 2=done)
// over edges (name -> reachable names). label and lineOf are used to format
// an error naming the cycle's re-entry point at its declaration line.
func detectNameCycle(path, label string, edges map[string][]string, lineOf func(string) int) error {
	state := make(map[string]int)
	var visit func(name string) error
	visit = func(name string) error {
		if state[name] == 1 {
			return errAt(path, lineOf(name), "%s cycle detected at %q", label, name)
		}
		if state[name] == 2 {
			return nil
		}
		state[name] = 1
		for _, next := range edges[name] {
			if err := visit(next); err != nil {
				return err
			}
		}
		state[name] = 2
		return nil
	}
	for name := range edges {
		if err := visit(name); err != nil {
			return err
		}
	}
	return nil
}
