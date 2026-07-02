//ff:func feature=tangl type=validator control=sequence
//ff:what TestBuildCheckingEdges — tests buildCheckingEdges for no-cases, empty-checking-skip, and non-empty-checking-append branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildCheckingEdges(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		edges := buildCheckingEdges(doc)
		if len(edges) != 0 {
			t.Fatalf("expected empty edges map, got %v", edges)
		}
	})

	t.Run("SkipsEmptyChecking", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: ""},
					},
				},
			},
		}
		edges := buildCheckingEdges(doc)
		if len(edges) != 0 {
			t.Fatalf("expected no edges for empty checking, got %v", edges)
		}
	})

	t.Run("AppendsMultiple", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Nodes: []ast.Node{
						{Name: "a", Checking: "y"},
						{Name: "b", Checking: ""},
						{Name: "c", Checking: "z"},
					},
				},
				{
					Name: "y",
					Nodes: []ast.Node{
						{Name: "d", Checking: "x"},
					},
				},
			},
		}
		edges := buildCheckingEdges(doc)
		want := map[string][]string{
			"x": {"y", "z"},
			"y": {"x"},
		}
		if !reflect.DeepEqual(edges, want) {
			t.Fatalf("expected edges=%v, got %v", want, edges)
		}
	})
}
