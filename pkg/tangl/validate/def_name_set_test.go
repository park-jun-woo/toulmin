//ff:func feature=tangl type=validator control=sequence
//ff:what TestDefNameSet — tests defNameSet for empty and populated Defs branches
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestDefNameSet(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		set := defNameSet(doc)
		if len(set) != 0 {
			t.Fatalf("expected empty set, got %v", set)
		}
	})

	t.Run("Populated", func(t *testing.T) {
		doc := &ast.Document{
			Defs: []ast.Definition{
				{Name: "a"},
				{Name: "b"},
			},
		}
		set := defNameSet(doc)
		if !set["a"] || !set["b"] || len(set) != 2 {
			t.Fatalf("expected set with a and b, got %v", set)
		}
	})
}
