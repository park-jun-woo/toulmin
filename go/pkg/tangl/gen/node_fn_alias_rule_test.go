//ff:func feature=tangl type=codegen control=sequence
//ff:what TestNodeFnAliasAndRule — alias checking targets the sibling Checking wrapper; bound host rules supply the function
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestNodeFnAliasAndRule(t *testing.T) {
	gc := &genContext{Rules: map[string]ast.InlineRule{
		"tries exhausted": {Name: "tries exhausted", Using: &ast.Ref{Alias: "ratchet", Name: "triesExhausted"}},
		"expr":            {Name: "expr", Cond: ast.Compare{Field: "x", Op: "is empty"}},
	}}
	if got, err := nodeFn(gc, ast.Node{Name: "n", Checking: &ast.Ref{Alias: "gate", Name: "can pass"}}); err != nil || got != "gate.CheckingCanPass" {
		t.Fatalf("alias checking = %q, %v", got, err)
	}
	if got, err := nodeFn(gc, ast.Node{Name: "tries exhausted"}); err != nil || got != "ratchet.triesExhausted" {
		t.Fatalf("bound host rule = %q, %v", got, err)
	}
	if got, err := nodeFn(gc, ast.Node{Name: "expr"}); err != nil || got != "expr" {
		t.Fatalf("expression rule = %q, %v", got, err)
	}
}
