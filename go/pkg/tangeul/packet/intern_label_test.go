//ff:func feature=tangl type=util control=sequence
//ff:what TestInternLabel — a new string appends and gets the next index; a repeated string reuses its index
package packet

import (
	"testing"
)

// TestInternLabel checks internLabel appends a new string and returns its
// index, that a later distinct string appends past it, and that a repeated
// string (found mid-scan) returns its existing index without growing the
// table.
func TestInternLabel(t *testing.T) {
	var labels []string

	if idx := internLabel(&labels, "a"); idx != 0 {
		t.Fatalf("first intern: got %d, want 0", idx)
	}
	if idx := internLabel(&labels, "b"); idx != 1 {
		t.Fatalf("second intern: got %d, want 1", idx)
	}
	if idx := internLabel(&labels, "a"); idx != 0 {
		t.Fatalf("repeat intern: got %d, want 0", idx)
	}
	if len(labels) != 2 {
		t.Fatalf("labels grew on repeat: got %v", labels)
	}
}
