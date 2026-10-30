//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileDocumentRuleErrors — sibling rule cycles and unresolvable host predicates are compile errors
package eval

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCompileDocumentRuleErrors(t *testing.T) {
	cyclic := &ast.Document{Path: "c.md", Rules: []ast.InlineRule{
		{Name: "a", Line: 1, Cond: ast.Call{Fn: ast.Ref{Name: "b"}}},
		{Name: "b", Line: 2, Cond: ast.Call{Fn: ast.Ref{Name: "a"}}},
	}}
	if fns, err := CompileDocument(cyclic, CompileOptions{}); err == nil || fns != nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("cycle: fns=%v err=%v", fns, err)
	}
	unresolved := &ast.Document{Rules: []ast.InlineRule{{Name: "h", Using: &ast.Ref{Alias: "nlp", Name: "spam"}}}}
	if fns, err := CompileDocument(unresolved, CompileOptions{Resolver: fakeResolver{}}); err == nil || fns != nil || !strings.Contains(err.Error(), "unresolved host predicate") {
		t.Fatalf("unresolved: fns=%v err=%v", fns, err)
	}
}
