//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what detectRuleCycle — 3-color DFS cycle detection over sibling-rule edges
package eval

import "fmt"

// detectRuleCycle performs a 3-color DFS (0=unvisited, 1=visiting, 2=done) over
// edges (rule name -> referenced sibling rule names), returning a formatted
// compile error naming the cycle's re-entry point at its declaration line.
func detectRuleCycle(path string, edges map[string][]string, lineOf func(string) int) error {
	state := map[string]int{}
	var visit func(name string) error
	visit = func(name string) error {
		if state[name] == 1 {
			return fmt.Errorf("%s:%d: rule cycle detected at %q", path, lineOf(name), name)
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
