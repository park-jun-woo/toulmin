//ff:func feature=tangl type=codegen control=sequence
//ff:what TestCollectUsedAliasesRules — host rule aliases are collected; bare host refs are not
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCollectUsedAliasesRules(t *testing.T) {
	used := collectUsedAliases(&ast.Document{Rules: []ast.InlineRule{
		{Name: "a", Using: &ast.Ref{Alias: "ratchet", Name: "f"}},
		{Name: "b", Using: &ast.Ref{Name: "local"}},
		{Name: "c", Cond: ast.Compare{Field: "x", Op: "is empty"}},
	}})
	if len(used) != 1 || !used["ratchet"] {
		t.Fatalf("used = %v, want only ratchet", used)
	}
}
