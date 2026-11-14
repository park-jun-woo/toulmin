//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestDetectCycle — tests DetectCycle for empty-edges, no-cycle, revisit-done-node, and cycle-detected branches
package toulmin

import (
	"strings"
	"testing"
)

func TestDetectCycle(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"EmptyEdges", func(t *testing.T) {
			if err := DetectCycle(map[string][]string{}); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		}},
		{"NoCycle", func(t *testing.T) {
			edges := map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {},
			}
			if err := DetectCycle(edges); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		}},
		{"RevisitsDoneNode", func(t *testing.T) {
			// Diamond: a->b, a->c, b->d, c->d. When d is visited a second time
			// (via c after being finished via b), state[d]==2 hits the early-return branch.
			edges := map[string][]string{
				"a": {"b", "c"},
				"b": {"d"},
				"c": {"d"},
				"d": {},
			}
			if err := DetectCycle(edges); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		}},
		{"DetectsCycle", func(t *testing.T) {
			edges := map[string][]string{
				"a": {"b"},
				"b": {"a"},
			}
			err := DetectCycle(edges)
			if err == nil {
				t.Fatal("expected cycle error, got nil")
			}
			if !strings.Contains(err.Error(), "cycle detected") {
				t.Errorf("expected cycle detected error, got %v", err)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
