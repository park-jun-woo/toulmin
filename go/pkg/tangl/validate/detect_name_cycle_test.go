//ff:func feature=tangl type=validator control=sequence
//ff:what TestDetectNameCycle — tests detectNameCycle for empty-edges, no-cycle, revisit-done-node, and cycle-detected branches
package validate

import (
	"strings"
	"testing"
)

func TestDetectNameCycle(t *testing.T) {
	t.Run("EmptyEdges", func(t *testing.T) {
		edges := map[string][]string{}
		err := detectNameCycle("doc.md", "run", edges, func(string) int { return 0 })
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("NoCycle", func(t *testing.T) {
		edges := map[string][]string{
			"a": {"b"},
			"b": {"c"},
			"c": {},
		}
		err := detectNameCycle("doc.md", "run", edges, func(string) int { return 0 })
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("RevisitsDoneNode", func(t *testing.T) {
		// Diamond: a->b, a->c, b->d, c->d. When d is visited a second time
		// (via c after being finished via b), state[d]==2 hits the early-return branch.
		edges := map[string][]string{
			"a": {"b", "c"},
			"b": {"d"},
			"c": {"d"},
			"d": {},
		}
		err := detectNameCycle("doc.md", "run", edges, func(string) int { return 0 })
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("DetectsCycle", func(t *testing.T) {
		edges := map[string][]string{
			"a": {"b"},
			"b": {"a"},
		}
		lines := map[string]int{"a": 1, "b": 2}
		err := detectNameCycle("doc.md", "run", edges, func(name string) int { return lines[name] })
		if err == nil {
			t.Fatal("expected cycle error, got nil")
		}
		if !strings.Contains(err.Error(), "run cycle detected") {
			t.Errorf("expected run cycle detected error, got %v", err)
		}
	})
}
