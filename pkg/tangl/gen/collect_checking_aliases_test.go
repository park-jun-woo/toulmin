//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectCheckingAliases — covers alias collection across alias, local, and no-checking nodes
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestCollectCheckingAliases(t *testing.T) {
	doc := &ast.Document{Cases: []ast.Case{{
		Nodes: []ast.Node{
			{Name: "n1", Checking: &ast.Ref{Alias: "b", Name: "cb"}}, // collected
			{Name: "n2"}, // no checking clause
			{Name: "n3", Checking: &ast.Ref{Name: "local"}},          // local, empty alias
			{Name: "n4", Checking: &ast.Ref{Alias: "b", Name: "cc"}}, // duplicate alias
		},
	}}}
	got := collectCheckingAliases(doc)
	if len(got) != 1 || !got["b"] {
		t.Fatalf("got %v, want only {b}", got)
	}
}
