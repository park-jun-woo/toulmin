//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildCases — tests buildCases for the success loop and the renderCase-error propagation branch
package gen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestBuildCases(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		doc := &ast.Document{
			Subject: "t",
			Cases: []ast.Case{
				{
					Name:  "one",
					Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}},
				},
				{
					Name:  "two",
					Nodes: []ast.Node{{Name: "n2", Role: ast.GeneralRule}},
				},
			},
		}
		gc := &genContext{Doc: doc, Defs: map[string]defInfo{}, CheckWrappers: map[string]string{}}
		var w strings.Builder
		if err := buildCases(&w, gc); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out := w.String()
		if !strings.Contains(out, "newOneGraph") || !strings.Contains(out, "newTwoGraph") {
			t.Errorf("expected both case functions rendered, got:\n%s", out)
		}
	})

	t.Run("renderCase error propagates", func(t *testing.T) {
		// A node with two "run" edges is rejected by runTarget, and buildCases
		// must return that error from the loop's err != nil branch.
		doc := &ast.Document{
			Subject: "t",
			Cases: []ast.Case{
				{
					Name:  "bad",
					Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}},
					Execs: []ast.Exec{
						{Kind: ast.RunExec, Case: "other1", Node: "n1"},
						{Kind: ast.RunExec, Case: "other2", Node: "n1"},
					},
				},
			},
		}
		gc := &genContext{Doc: doc, Defs: map[string]defInfo{}, CheckWrappers: map[string]string{}}
		var w strings.Builder
		if err := buildCases(&w, gc); err == nil {
			t.Fatal("expected error for node with multiple run edges")
		}
	})
}
