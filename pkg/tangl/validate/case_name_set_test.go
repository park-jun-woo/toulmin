//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseNameSet — tests caseNameSet for empty-cases and multiple-cases branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCaseNameSet(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		set := caseNameSet(doc)
		if len(set) != 0 {
			t.Fatalf("expected empty set, got %v", set)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x"},
				{Name: "y"},
			},
		}
		set := caseNameSet(doc)
		want := map[string]bool{"x": true, "y": true}
		if !reflect.DeepEqual(set, want) {
			t.Fatalf("expected set=%v, got %v", want, set)
		}
	})
}
