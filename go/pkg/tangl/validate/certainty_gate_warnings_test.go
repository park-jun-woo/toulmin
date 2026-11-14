//ff:func feature=tangl type=validator control=sequence
//ff:what TestCertaintyGateWarnings — never-passing and always-passing gates warn; ungated and in-range gates stay silent
package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestCertaintyGateWarnings(t *testing.T) {
	doc := &ast.Document{Path: "g.md", Cases: []ast.Case{{
		Name: "c",
		Execs: []ast.Exec{
			{Kind: ast.DoExec, Node: "n"},
			{Kind: ast.DoExec, Node: "n", Certainty: &ast.Certainty{Op: "at least", Percent: 75}},
			{Kind: ast.DoExec, Node: "n", Certainty: &ast.Certainty{Op: "at most", Percent: 50}, Line: 4},
			{Kind: ast.DoExec, Node: "n", Certainty: &ast.Certainty{Op: "at least", Percent: 40}, Line: 5},
		},
	}}}
	w := certaintyGateWarnings(doc)
	if len(w) != 2 {
		t.Fatalf("got %v, want 2 warnings", w)
	}
	if !strings.HasPrefix(w[0], "g.md:4: warning:") || !strings.Contains(w[0], "'if at most 50% certain' never passes") {
		t.Errorf("w[0] = %q", w[0])
	}
	if !strings.HasPrefix(w[1], "g.md:5: warning:") || !strings.Contains(w[1], "'if at least 40% certain' always passes") {
		t.Errorf("w[1] = %q", w[1])
	}
}
