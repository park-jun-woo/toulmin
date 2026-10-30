//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateWithOptionsBranches — import, rule and case errors abort; a case pulls in the case helper
package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestGenerateWithOptionsBranches(t *testing.T) {
	importErr := &ast.Document{
		Subject: "t",
		Path:    filepath.Join(t.TempDir(), "doc.geul"),
		Sees:    []ast.See{{Alias: "gate", Package: "gate"}},
		Cases:   []ast.Case{{Name: "c", Nodes: []ast.Node{{Name: "n", Role: ast.GeneralRule, Checking: &ast.Ref{Alias: "gate", Name: "x"}}}}},
	}
	if out, err := GenerateWithOptions(importErr, GenerateOptions{}); err == nil || out != "" {
		t.Fatalf("import error: out=%q err=%v", out, err)
	}

	ruleErr := &ast.Document{Subject: "t", Rules: []ast.InlineRule{{Name: "bad"}}}
	if out, err := GenerateWithOptions(ruleErr, GenerateOptions{}); err == nil || out != "" {
		t.Fatalf("rule error: out=%q err=%v", out, err)
	}

	caseErr := &ast.Document{Subject: "t", Cases: []ast.Case{{
		Name:  "bad case",
		Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}},
		Execs: []ast.Exec{{Kind: ast.RunExec, Case: "a", Node: "n1"}, {Kind: ast.RunExec, Case: "b", Node: "n1"}},
	}}}
	if out, err := GenerateWithOptions(caseErr, GenerateOptions{}); err == nil || out != "" {
		t.Fatalf("case error: out=%q err=%v", out, err)
	}

	ok := &ast.Document{Subject: "t", Cases: []ast.Case{{Name: "gate", Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}}}}}
	out, err := GenerateWithOptions(ok, GenerateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "func newGateGraph() *toulmin.Graph") {
		t.Fatalf("case graph builder missing:\n%s", out)
	}
	var helper strings.Builder
	writeCaseHelper(&helper)
	if first := strings.SplitN(strings.TrimSpace(helper.String()), "\n", 2)[0]; !strings.Contains(out, first) {
		t.Fatalf("case helper %q missing from output", first)
	}
}
