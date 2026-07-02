//ff:func feature=tangl type=analyzer control=iteration dimension=1
//ff:what Closure — static do/undo transitive closure reachable from an endpoint's run cases
package effects

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// Closure computes the static do/undo effect summary reachable from
// endpoint's `run` cases, following Cases' `run <case> when <node>` edges in
// document declaration order. `checking` references are excluded from the
// closure (they compile to Evaluate and stay pure). A check-only endpoint
// (no `run` cases) always yields an empty, non-nil slice. An endpoint name
// unknown in doc is an error.
func Closure(doc *ast.Document, endpoint string) ([]Entry, error) {
	ep := findEndpoint(doc, endpoint)
	if ep == nil {
		return nil, fmt.Errorf("tangl: unknown endpoint %q", endpoint)
	}
	entries := []Entry{}
	if len(ep.Runs) == 0 {
		return entries, nil
	}
	state := make(map[string]int) // 0=unvisited, 1=visiting, 2=done
	var visit func(name string) error
	visit = func(name string) error {
		if state[name] == 1 || state[name] == 2 {
			return nil
		}
		state[name] = 1
		c := findCase(doc, name)
		if c == nil {
			return fmt.Errorf("tangl: endpoint %q runs unknown case %q", endpoint, name)
		}
		for _, e := range c.Execs {
			switch e.Kind {
			case ast.DoExec:
				entries = append(entries, Entry{Kind: "do", Func: *e.Func, Once: e.Once, Case: c.Name, Node: e.Node})
			case ast.UndoExec:
				entries = append(entries, Entry{Kind: "undo", Func: *e.Func, Case: c.Name, Node: e.Node})
			case ast.RunExec:
				if err := visit(e.Case); err != nil {
					return err
				}
			}
		}
		state[name] = 2
		return nil
	}
	for _, name := range ep.Runs {
		if err := visit(name); err != nil {
			return nil, err
		}
	}
	return entries, nil
}
