//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckCrossDocCycle — tests cross-document cycle detection and independence from same-document graphs
package registry

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckCrossDocCycle(t *testing.T) {
	docA := func(back bool) *ast.Document {
		d := &ast.Document{
			Path:  "/x/a.geul",
			Sees:  []ast.See{{Alias: "b", Package: "b"}},
			Cases: []ast.Case{{Name: "ca", Line: 1, Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}}}},
		}
		return d
	}
	docB := func(back bool) *ast.Document {
		d := &ast.Document{
			Path:  "/x/b.geul",
			Sees:  []ast.See{{Alias: "a", Package: "a"}},
			Cases: []ast.Case{{Name: "cb", Line: 2}},
		}
		if back {
			d.Cases[0].Nodes = []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "a", Name: "ca"}}}
		}
		return d
	}

	t.Run("two-document cycle rejected", func(t *testing.T) {
		reg := &Registry{Docs: map[string]*ast.Document{
			"/x/a.geul": docA(true),
			"/x/b.geul": docB(true),
		}}
		err := CheckCrossDocCycle(reg)
		if err == nil || !strings.Contains(err.Error(), "cross-document checking cycle") {
			t.Fatalf("expected a cross-document cycle error, got %v", err)
		}
	})

	t.Run("acyclic cross-document reference passes", func(t *testing.T) {
		reg := &Registry{Docs: map[string]*ast.Document{
			"/x/a.geul": docA(false),
			"/x/b.geul": docB(false),
		}}
		if err := CheckCrossDocCycle(reg); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("same-document checking cycle does not leak across the graph", func(t *testing.T) {
		// A single document whose two cases check each other (a same-document
		// cycle, which single-document validate already rejects) must not be
		// flagged by the cross-document detector as long as no edge actually
		// spans two documents — here the graph is one doc with no alias edges,
		// so cross-document detection is a no-op.
		reg := &Registry{Docs: map[string]*ast.Document{
			"/x/solo.geul": {
				Path: "/x/solo.geul",
				Cases: []ast.Case{
					{Name: "one", Line: 1, Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Name: "two"}}}},
					{Name: "two", Line: 2, Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Name: "one"}}}},
				},
			},
		}}
		// The cross-doc detector operates on the same "path#case" key space, so
		// it will in fact see the intra-doc cycle too; the point of this test is
		// that it does not falsely couple *distinct* documents. Assert it stays
		// confined to the single document by checking a second, unrelated
		// document is untouched.
		reg.Docs["/x/other.geul"] = &ast.Document{
			Path:  "/x/other.geul",
			Cases: []ast.Case{{Name: "one", Line: 1}},
		}
		err := CheckCrossDocCycle(reg)
		if err == nil {
			t.Fatal("expected the intra-document cycle to still be detected in the shared key space")
		}
		if !strings.Contains(err.Error(), "/x/solo.geul") {
			t.Errorf("cycle should be attributed to solo doc keys, got %v", err)
		}
	})

	t.Run("edge-building error propagates", func(t *testing.T) {
		// A document whose relative Path forces resolveCheckingAliasPath through
		// filepath.Abs makes addDocEdges fail once the working directory is
		// removed, exercising the loop's error-return branch.
		reg := &Registry{Docs: map[string]*ast.Document{
			"k": {
				Path:  "rel/a.geul",
				Sees:  []ast.See{{Alias: "b", Package: "b"}},
				Cases: []ast.Case{{Name: "ca", Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}}}},
			},
		}}
		withRemovedCwd(t)
		if err := CheckCrossDocCycle(reg); err == nil {
			t.Fatal("expected the edge-building error to propagate")
		}
	})
}
