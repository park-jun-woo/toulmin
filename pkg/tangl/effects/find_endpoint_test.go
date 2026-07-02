//ff:func feature=tangl type=analyzer control=sequence
//ff:what TestFindEndpoint — tests findEndpoint for found, not-found, and empty-provides branches
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestFindEndpoint(t *testing.T) {
	doc := &ast.Document{
		Provides: []ast.Endpoint{
			{Name: "a"},
			{Name: "b"},
		},
	}

	if e := findEndpoint(doc, "b"); e == nil || e.Name != "b" {
		t.Errorf("expected to find endpoint %q, got %+v", "b", e)
	}

	if e := findEndpoint(doc, "missing"); e != nil {
		t.Errorf("expected nil for missing endpoint, got %+v", e)
	}

	emptyDoc := &ast.Document{Provides: []ast.Endpoint{}}
	if e := findEndpoint(emptyDoc, "anything"); e != nil {
		t.Errorf("expected nil for empty provides, got %+v", e)
	}
}
