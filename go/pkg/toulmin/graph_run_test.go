//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphRun — tests Graph.Run for run-cycle error, resolveOption error, and success branches
package toulmin

import (
	"strings"
	"testing"
)

func TestGraphRun(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"CycleError", func(t *testing.T) {
			f1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			ga := NewGraph("A")
			gb := NewGraph("B")
			ga.Rule(f1).Run(gb)
			gb.Rule(f1).Run(ga)

			_, _, err := ga.Run(NewContext())
			if err == nil {
				t.Fatal("expected run cycle error")
			}
			if !strings.Contains(err.Error(), "cycle") {
				t.Errorf("expected cycle error, got %v", err)
			}
		}},
		{"ResolveOptionError", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			_, _, err := g.Run(NewContext(), EvalOption{Method: Recursive})
			if err == nil {
				t.Fatal("expected resolveOption error for Recursive method")
			}
			if !strings.Contains(err.Error(), "not yet implemented") {
				t.Errorf("expected not yet implemented error, got %v", err)
			}
		}},
		{"Success", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			results, trace, err := g.Run(NewContext())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 result, got %d", len(results))
			}
			_ = trace
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
