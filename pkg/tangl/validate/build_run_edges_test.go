//ff:func feature=tangl type=validator control=sequence
//ff:what TestBuildRunEdges — tests buildRunEdges for no-cases, non-run-exec-skip, and run-exec-append branches
package validate

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildRunEdges(t *testing.T) {
	t.Run("NoCases", func(t *testing.T) {
		doc := &ast.Document{}
		edges := buildRunEdges(doc)
		if len(edges) != 0 {
			t.Fatalf("expected empty edges map, got %v", edges)
		}
	})

	t.Run("SkipsNonRunExec", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Kind: ast.DoExec, Case: "y"},
					},
				},
			},
		}
		edges := buildRunEdges(doc)
		if len(edges) != 0 {
			t.Fatalf("expected no edges for non-run exec, got %v", edges)
		}
	})

	t.Run("AppendsMultiple", func(t *testing.T) {
		doc := &ast.Document{
			Cases: []ast.Case{
				{
					Name: "x",
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "y"},
						{Kind: ast.DoExec, Case: "ignored"},
						{Kind: ast.RunExec, Case: "z"},
					},
				},
				{
					Name: "y",
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "x"},
					},
				},
			},
		}
		edges := buildRunEdges(doc)
		want := map[string][]string{
			"x": {"y", "z"},
			"y": {"x"},
		}
		if !reflect.DeepEqual(edges, want) {
			t.Fatalf("expected edges=%v, got %v", want, edges)
		}
	})
}
