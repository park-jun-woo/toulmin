//ff:func feature=tangl type=engine control=sequence
//ff:what TestCompileDocument — materialization and rule-Spec errors abort; derived Specs reach the compiled rule
package eval

import (
	"errors"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

func TestCompileDocument(t *testing.T) {
	fileDoc := &ast.Document{Path: "/p.md", Defs: []ast.Definition{{Name: "f", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "x.json", Line: 2}}}
	if _, err := CompileDocument(fileDoc, CompileOptions{}); err == nil || !strings.Contains(err.Error(), "loader registry is required") {
		t.Fatalf("missing loader registry: err = %v", err)
	}

	specDoc := &ast.Document{
		Defs:  []ast.Definition{{Name: "limit", Kind: ast.ConstDef, Value: "3", SpecRef: &ast.Ref{Name: "NewLimit"}}},
		Rules: []ast.InlineRule{{Name: "r", Cond: ast.Compare{Lhs: ast.SpecRef{Steps: []string{"limit"}}, Op: "is not empty"}}},
	}
	if _, err := CompileDocument(specDoc, CompileOptions{Resolver: failingDefinitionResolver{err: errors.New("boom")}}); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("rule Spec error: err = %v", err)
	}

	fns, err := CompileDocument(specDoc, CompileOptions{Resolver: definitionSpecResolver{}})
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := fns["r"](toulmin.NewContext(), nil); !got {
		t.Fatal("derived Spec limit must be visible to the rule")
	}
}
