//ff:func feature=tangl type=validator control=sequence
//ff:what TestDetectNameCycle_Export — tests the exported DetectNameCycle wrapper delegates to the cycle detector
package validate

import "testing"

func TestDetectNameCycle_Export(t *testing.T) {
	lineOf := func(string) int { return 0 }

	t.Run("cycle detected", func(t *testing.T) {
		edges := map[string][]string{"a": {"b"}, "b": {"a"}}
		if err := DetectNameCycle("p", "test", edges, lineOf); err == nil {
			t.Fatal("expected a cycle error")
		}
	})

	t.Run("acyclic graph passes", func(t *testing.T) {
		edges := map[string][]string{"a": {"b"}, "b": {"c"}}
		if err := DetectNameCycle("p", "test", edges, lineOf); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
