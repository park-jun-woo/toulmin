//ff:func feature=tangl type=validator control=sequence
//ff:what TestAddCaseEdges — covers edge append, no-checking skip, and resolve-error propagation
package registry

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestAddCaseEdges(t *testing.T) {
	t.Run("appends resolved edge and skips no-checking node", func(t *testing.T) {
		d := &ast.Document{Path: "/repo/a.geul", Sees: []ast.See{{Alias: "b", Package: "b"}}}
		c := ast.Case{Name: "ca", Nodes: []ast.Node{
			{Name: "n1", Checking: &ast.Ref{Alias: "b", Name: "cb"}},
			{Name: "n2"}, // no checking clause -> contributes no edge
		}}
		edges := map[string][]string{}
		if err := addCaseEdges(edges, "/repo/a.geul", d, c, "/repo/a.geul#ca"); err != nil {
			t.Fatal(err)
		}
		got := edges["/repo/a.geul#ca"]
		if len(got) != 1 || got[0] != "/repo/b.geul#cb" {
			t.Fatalf("got %v, want one edge to /repo/b.geul#cb", got)
		}
	})

	t.Run("checkingEdgeTarget error propagates", func(t *testing.T) {
		rel := &ast.Document{Path: "rel/a.geul", Sees: []ast.See{{Alias: "b", Package: "b"}}}
		c := ast.Case{Name: "ca", Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}}}
		withRemovedCwd(t)
		if err := addCaseEdges(map[string][]string{}, "rel/a.geul", rel, c, "k"); err == nil {
			t.Fatal("expected the resolve error to propagate")
		}
	})
}
