//ff:func feature=tangl type=analyzer control=sequence
//ff:what TestFindCase — tests findCase for found, not-found, and empty-cases branches
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestFindCase(t *testing.T) {
	doc := &ast.Document{
		Cases: []ast.Case{
			{Name: "a"},
			{Name: "b"},
		},
	}

	if c := findCase(doc, "b"); c == nil || c.Name != "b" {
		t.Errorf("expected to find case %q, got %+v", "b", c)
	}

	if c := findCase(doc, "missing"); c != nil {
		t.Errorf("expected nil for missing case, got %+v", c)
	}

	emptyDoc := &ast.Document{Cases: []ast.Case{}}
	if c := findCase(emptyDoc, "anything"); c != nil {
		t.Errorf("expected nil for empty cases, got %+v", c)
	}
}
