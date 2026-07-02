//ff:func feature=tangl type=validator control=sequence
//ff:what TestCaseLineIndex — tests caseLineIndex for empty-cases and multiple-cases branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCaseLineIndex(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		doc := &ast.Document{}
		idx := caseLineIndex(doc)
		if len(idx) != 0 {
			t.Fatalf("expected empty index, got %v", idx)
		}
	})

	t.Run("Multiple", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{Name: "x", Line: 3},
				{Name: "y", Line: 7},
			},
		}
		idx := caseLineIndex(doc)
		want := map[string]int{"x": 3, "y": 7}
		if !reflect.DeepEqual(idx, want) {
			t.Fatalf("expected idx=%v, got %v", want, idx)
		}
	})
}
