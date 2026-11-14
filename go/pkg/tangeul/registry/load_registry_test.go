//ff:func feature=tangl type=util control=sequence
//ff:what TestLoadRegistry — tests recursive load, file-load cycle rejection, diamond caching, missing-See, and skipped-node paths via an injected loader
package registry

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestLoadRegistry(t *testing.T) {
	t.Run("two-level reference loads both", func(t *testing.T) {
		docs := map[string]*ast.Document{
			"/base/a.geul": {
				Sees:  []ast.See{{Alias: "b", Package: "b"}},
				Cases: []ast.Case{aliasCheckNode("ca", "b", "cb")},
			},
			"/base/b.geul": {Cases: []ast.Case{{Name: "cb"}}},
		}
		reg, err := loadRegistry("/base/a.geul", fakeLoader(docs, map[string]int{}))
		if err != nil {
			t.Fatal(err)
		}
		if len(reg.Docs) != 2 || reg.Docs["/base/a.geul"] == nil || reg.Docs["/base/b.geul"] == nil {
			t.Errorf("expected both canonical paths, got %d docs", len(reg.Docs))
		}
	})

	t.Run("file-load cycle rejected", func(t *testing.T) {
		docs := map[string]*ast.Document{
			"/base/a.geul": {
				Sees:  []ast.See{{Alias: "b", Package: "b"}},
				Cases: []ast.Case{aliasCheckNode("ca", "b", "cb")},
			},
			"/base/b.geul": {
				Sees:  []ast.See{{Alias: "a", Package: "a"}},
				Cases: []ast.Case{aliasCheckNode("cb", "a", "ca")},
			},
		}
		_, err := loadRegistry("/base/a.geul", fakeLoader(docs, map[string]int{}))
		if err == nil || !strings.Contains(err.Error(), "load cycle") {
			t.Fatalf("expected a load-cycle error, got %v", err)
		}
	})

	t.Run("diamond parses shared leaf once", func(t *testing.T) {
		docs := map[string]*ast.Document{
			"/base/a.geul": {
				Sees: []ast.See{{Alias: "b", Package: "b"}, {Alias: "c", Package: "c"}},
				Cases: []ast.Case{{Name: "ca", Nodes: []ast.Node{
					{Name: "nb", Checking: &ast.Ref{Alias: "b", Name: "cb"}},
					{Name: "nc", Checking: &ast.Ref{Alias: "c", Name: "cc"}},
				}}},
			},
			"/base/b.geul": {Sees: []ast.See{{Alias: "d", Package: "d"}}, Cases: []ast.Case{aliasCheckNode("cb", "d", "cd")}},
			"/base/c.geul": {Sees: []ast.See{{Alias: "d", Package: "d"}}, Cases: []ast.Case{aliasCheckNode("cc", "d", "cd")}},
			"/base/d.geul": {Cases: []ast.Case{{Name: "cd"}}},
		}
		counts := map[string]int{}
		reg, err := loadRegistry("/base/a.geul", fakeLoader(docs, counts))
		if err != nil {
			t.Fatal(err)
		}
		if len(reg.Docs) != 4 {
			t.Fatalf("expected 4 docs, got %d", len(reg.Docs))
		}
		if counts["/base/d.geul"] != 1 {
			t.Errorf("expected shared leaf parsed once, got %d", counts["/base/d.geul"])
		}
	})

	t.Run("checking alias without See declaration errors", func(t *testing.T) {
		docs := map[string]*ast.Document{
			"/base/a.geul": {Cases: []ast.Case{aliasCheckNode("ca", "missing", "cb")}},
		}
		_, err := loadRegistry("/base/a.geul", fakeLoader(docs, map[string]int{}))
		if err == nil || !strings.Contains(err.Error(), "no matching See") {
			t.Fatalf("expected a missing-See error, got %v", err)
		}
	})

	t.Run("node without a checking alias is skipped", func(t *testing.T) {
		// A plain node (no checking clause) exercises the continue branch, and
		// a document with no alias references loads as a single-doc registry.
		docs := map[string]*ast.Document{
			"/base/a.geul": {Cases: []ast.Case{{Name: "ca", Nodes: []ast.Node{{Name: "n"}}}}},
		}
		reg, err := loadRegistry("/base/a.geul", fakeLoader(docs, map[string]int{}))
		if err != nil {
			t.Fatal(err)
		}
		if len(reg.Docs) != 1 {
			t.Fatalf("expected 1 doc, got %d", len(reg.Docs))
		}
	})
}
