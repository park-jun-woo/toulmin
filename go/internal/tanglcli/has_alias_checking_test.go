//ff:func feature=cli type=command control=sequence
//ff:what TestHasAliasChecking — covers alias-checking detection across nil, empty-alias, and matching nodes
package tanglcli

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestHasAliasChecking(t *testing.T) {
	t.Run("node with alias-qualified checking returns true", func(t *testing.T) {
		d := &ast.Document{Cases: []ast.Case{{
			Name:  "c",
			Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "b", Name: "cb"}}},
		}}}
		if !hasAliasChecking(d) {
			t.Fatal("expected true for an alias-qualified checking clause")
		}
	})

	t.Run("no alias-qualified checking returns false", func(t *testing.T) {
		// One node has no checking clause; another checks a local (empty-alias)
		// target. Neither is alias-qualified, so the scan falls through to false.
		d := &ast.Document{Cases: []ast.Case{{
			Name: "c",
			Nodes: []ast.Node{
				{Name: "n1"},
				{Name: "n2", Checking: &ast.Ref{Name: "local"}},
			},
		}}}
		if hasAliasChecking(d) {
			t.Fatal("expected false when no checking clause is alias-qualified")
		}
	})
}
