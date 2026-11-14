//ff:func feature=approve type=engine control=sequence
//ff:what TestNewFlow — NewFlow builds an empty Flow with the given name
package approve

import "testing"

// TestNewFlow covers the single branch of NewFlow: it always returns a
// non-nil *Flow with the given name and no steps.
func TestNewFlow(t *testing.T) {
	f := NewFlow("expense")
	if f == nil {
		t.Fatal("NewFlow() returned nil")
	}
	if f.name != "expense" {
		t.Errorf("name = %q, want %q", f.name, "expense")
	}
	if len(f.steps) != 0 {
		t.Errorf("len(steps) = %d, want 0", len(f.steps))
	}
}
