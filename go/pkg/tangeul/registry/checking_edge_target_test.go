//ff:func feature=tangl type=validator control=sequence
//ff:what TestCheckingEdgeTarget — covers no-checking, local, alias, missing-See, and resolve-error paths
package registry

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCheckingEdgeTarget(t *testing.T) {
	d := &ast.Document{
		Path: "/repo/a.geul",
		Sees: []ast.See{{Alias: "b", Package: "b"}},
	}

	t.Run("no checking clause yields empty key", func(t *testing.T) {
		got, err := checkingEdgeTarget(d, "/repo/a.geul", ast.Node{Name: "n"})
		if err != nil || got != "" {
			t.Fatalf("got %q err=%v, want empty", got, err)
		}
	})

	t.Run("local target keeps docPath", func(t *testing.T) {
		n := ast.Node{Name: "n", Checking: &ast.Ref{Name: "cb"}}
		got, err := checkingEdgeTarget(d, "/repo/a.geul", n)
		if err != nil || got != "/repo/a.geul#cb" {
			t.Fatalf("got %q err=%v", got, err)
		}
	})

	t.Run("alias target resolves to sibling document", func(t *testing.T) {
		n := ast.Node{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}
		got, err := checkingEdgeTarget(d, "/repo/a.geul", n)
		if err != nil || got != "/repo/b.geul#cb" {
			t.Fatalf("got %q err=%v", got, err)
		}
	})

	t.Run("alias without a See declaration yields empty key", func(t *testing.T) {
		n := ast.Node{Name: "n", Checking: &ast.Ref{Alias: "unknown", Name: "cb"}}
		got, err := checkingEdgeTarget(d, "/repo/a.geul", n)
		if err != nil || got != "" {
			t.Fatalf("got %q err=%v, want empty", got, err)
		}
	})

	t.Run("resolve error propagates", func(t *testing.T) {
		// A relative document path forces resolveCheckingAliasPath through
		// filepath.Abs, which fails when the working directory has been removed.
		rel := &ast.Document{
			Path: "rel/a.geul",
			Sees: []ast.See{{Alias: "b", Package: "b"}},
		}
		withRemovedCwd(t)
		n := ast.Node{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}
		if _, err := checkingEdgeTarget(rel, "rel/a.geul", n); err == nil || !strings.Contains(err.Error(), "getwd") {
			t.Fatalf("expected a getwd error, got %v", err)
		}
	})
}
