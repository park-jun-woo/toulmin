//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerate_ErrorBranches — covers Generate's buildRules-error and buildCases-error propagation branches
package gen

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
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
}
