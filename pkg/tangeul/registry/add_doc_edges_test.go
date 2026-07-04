//ff:func feature=tangl type=validator control=sequence
//ff:what TestAddDocEdges — covers per-case line recording, edges, and error propagation
package registry

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestAddDocEdges(t *testing.T) {
	t.Run("records line and edges per case", func(t *testing.T) {
		d := &ast.Document{
			Path: "/repo/a.tangeul",
			Sees: []ast.See{{Alias: "b", Package: "b"}},
			Cases: []ast.Case{{
				Name:  "ca",
				Line:  7,
				Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}},
			}},
		}
		edges := map[string][]string{}
		lines := map[string]int{}
		if err := addDocEdges(edges, lines, "/repo/a.tangeul", d); err != nil {
			t.Fatal(err)
		}
		if lines["/repo/a.tangeul#ca"] != 7 {
			t.Errorf("line = %d, want 7", lines["/repo/a.tangeul#ca"])
		}
		if got := edges["/repo/a.tangeul#ca"]; len(got) != 1 || got[0] != "/repo/b.tangeul#cb" {
			t.Errorf("edges = %v", got)
		}
	})

	t.Run("addCaseEdges error propagates", func(t *testing.T) {
		rel := &ast.Document{
			Path: "rel/a.tangeul",
			Sees: []ast.See{{Alias: "b", Package: "b"}},
			Cases: []ast.Case{{
				Name:  "ca",
				Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}},
			}},
		}
		withRemovedCwd(t)
		if err := addDocEdges(map[string][]string{}, map[string]int{}, "rel/a.tangeul", rel); err == nil {
			t.Fatal("expected the resolve error to propagate")
		}
	})
}
