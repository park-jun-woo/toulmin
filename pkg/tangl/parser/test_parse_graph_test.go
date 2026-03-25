//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseGraph — test graph declaration and nested list parsing
package parser

import "testing"

// TestParseGraph tests graph declaration parsing with nested bindings.
func TestParseGraph(t *testing.T) {
	gd, err := parseGraph(`access is a graph "api:access"`, 5)
	if err != nil {
		t.Fatalf("parseGraph failed: %v", err)
	}
	if gd.Name != "access" {
		t.Errorf("expected name 'access', got %q", gd.Name)
	}
	if gd.ID != "api:access" {
		t.Errorf("expected id 'api:access', got %q", gd.ID)
	}

	input := `## tangl:Graph
- access is a graph "api:access"
  - auth is a rule using isAuthenticated
  - blocked is a counter using checkIP
  - blocked attacks auth
`
	f, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}
	if len(f.Graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(f.Graphs))
	}
	if len(f.Bindings) != 2 {
		t.Fatalf("expected 2 bindings, got %d", len(f.Bindings))
	}
	if f.Bindings[0].Graph != "access" {
		t.Errorf("expected binding graph 'access', got %q", f.Bindings[0].Graph)
	}
	if f.Bindings[1].Graph != "access" {
		t.Errorf("expected binding graph 'access', got %q", f.Bindings[1].Graph)
	}
	if len(f.Attacks) != 1 {
		t.Fatalf("expected 1 attack, got %d", len(f.Attacks))
	}
	if f.Attacks[0].Graph != "access" {
		t.Errorf("expected attack graph 'access', got %q", f.Attacks[0].Graph)
	}
}
