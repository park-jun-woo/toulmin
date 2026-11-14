//ff:func feature=engine type=engine control=sequence
//ff:what TestNewGraph — tests NewGraph returns a usable empty Graph
package toulmin

import "testing"

func TestNewGraph(t *testing.T) {
	g := NewGraph("mygraph")
	if g == nil {
		t.Fatalf("expected non-nil graph")
	}
	if g.name != "mygraph" {
		t.Fatalf("expected name %q, got %q", "mygraph", g.name)
	}
	if g.roles == nil {
		t.Fatalf("expected roles map initialized")
	}
	if len(g.roles) != 0 {
		t.Fatalf("expected empty roles map, got %v", g.roles)
	}
}
