//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerate_ErrorBranches — covers Generate's buildRules-error and buildCases-error propagation branches
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestGenerate_ErrorBranches(t *testing.T) {
	// verifies that an unrenderable tangl:Rules condition surfaces as an
	// error from Generate (the buildRules err != nil branch).
	t.Run("buildRules error propagates", func(t *testing.T) {
		doc := &ast.Document{
			Subject: "t",
			Rules: []ast.InlineRule{
				{Name: "bad rule", Cond: nil},
			},
		}
		if _, err := Generate(doc); err == nil {
			t.Fatal("expected error for unsupported rule condition")
		}
	})

	// verifies that a case node with more than one "run" edge surfaces as
	// an error from Generate (the buildCases err != nil branch).
	t.Run("buildCases error propagates", func(t *testing.T) {
		doc := &ast.Document{
			Subject: "t",
			Cases: []ast.Case{
				{
					Name:  "bad case",
					Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}},
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "other1", Node: "n1"},
						{Kind: ast.RunExec, Case: "other2", Node: "n1"},
					},
				},
			},
		}
		if _, err := Generate(doc); err == nil {
			t.Fatal("expected error for node with multiple run edges")
		}
	})

	// verifies that an alias-qualified checking target whose sibling document
	// cannot be located (no enclosing go.mod for the derived import path)
	// surfaces as an error from Generate (buildAliasImports err != nil branch)
	// with no partial output. The success path is covered by the cross-document
	// compile test.
	t.Run("alias-qualified checking without a module errors", func(t *testing.T) {
		doc := &ast.Document{
			Subject: "t",
			Path:    "/nonexistent-registry-root-xyz/doc.geul",
			Sees: []ast.See{
				{Alias: "gate", Package: "some/pkg"},
			},
			Cases: []ast.Case{
				{
					Name: "outer",
					Nodes: []ast.Node{
						{Name: "n1", Role: ast.GeneralRule, Checking: &ast.Ref{Alias: "gate", Name: "can pass"}},
					},
				},
			},
		}
		out, err := Generate(doc)
		if err == nil {
			t.Fatal("expected error deriving the sibling package import path")
		}
		if !strings.Contains(err.Error(), "go.mod") {
			t.Errorf("expected a go.mod-lookup error, got %v", err)
		}
		if out != "" {
			t.Errorf("expected no partial output, got %q", out)
		}
	})
}
