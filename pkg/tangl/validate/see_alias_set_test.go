//ff:func feature=tangl type=validator control=sequence
//ff:what TestSeeAliasSet — tests seeAliasSet for empty and populated Sees branches
package validate

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestSeeAliasSet(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		set := seeAliasSet(doc)
		if len(set) != 0 {
			t.Fatalf("expected empty set, got %v", set)
		}
	})

	t.Run("Populated", func(t *testing.T) {
		doc := &ast.Document{
			Sees: []ast.See{
				{Alias: "a"},
				{Alias: "b"},
			},
		}
		set := seeAliasSet(doc)
		if !set["a"] || !set["b"] || len(set) != 2 {
			t.Fatalf("expected set with a and b, got %v", set)
		}
	})
}
