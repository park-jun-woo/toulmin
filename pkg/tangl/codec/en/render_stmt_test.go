//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderStmt — RenderStmt covers the ResolveStmt error path and every AST element case
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderStmt directly exercises RenderStmt's own branches (independent of
// the round-trip corpus tests, which exercise the whole codec together): the
// ResolveStmt error short-circuit, the bare-string Subject and run/check forms,
// and the eleven typed-element switch cases.
func TestRenderStmt(t *testing.T) {
	fn := &ast.Ref{Alias: "p", Name: "check"}
	d := &doc.Doc{AST: &ast.Document{
		Subject: "demo",
		Sees:    []ast.See{{Alias: "p", Package: "pkg/p"}},
		Defs: []ast.Definition{
			{Name: "threshold", Kind: ast.ConstDef, Value: "10"},
			{Name: "applicant", Kind: ast.StructDef, Fields: []ast.Field{{Name: "income", Type: "Currency"}}},
		},
		Rules: []ast.InlineRule{{Name: "r", Cond: ast.Compare{Field: "f", Op: "equals", Value: "1"}}},
		Cases: []ast.Case{{
			Name:     "can go",
			Requires: []ast.Require{{Field: "user"}},
			Nodes:    []ast.Node{{Name: "ok", Role: ast.GeneralRule}},
			Attacks:  []ast.Attack{{Target: "ok", Attacker: "bad"}},
			Execs:    []ast.Exec{{Kind: ast.DoExec, Func: fn, Node: "ok"}},
		}},
		Provides: []ast.Endpoint{{
			Name:     "entry",
			Requires: []ast.Require{{Field: "user"}},
			Runs:     []string{"can go"},
			Checks:   []string{"can sub"},
		}},
		Internals: []ast.Internal{{
			Kind:   ast.OnEvent,
			Event:  "start",
			Runs:   []string{"can go"},
			Checks: []string{"can sub"},
		}},
	}}

	// ResolveStmt error: Subject Path must be nil, so this fails before the switch.
	if _, err := (ENCodec{}).RenderStmt(d, doc.StmtRef{Sec: codebook.SectionSubject, Path: []int{0}}); err == nil {
		t.Error("RenderStmt accepted a malformed Subject ref")
	}

	cases := []struct {
		name string
		ref  doc.StmtRef
		want string
	}{
		{"subject string", doc.StmtRef{Sec: codebook.SectionSubject, Path: nil}, "this document is `demo`"},
		{"see", doc.StmtRef{Sec: codebook.SectionSee, Path: []int{0}}, "see `p` from `pkg/p`"},
		{"const definition", doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0}}, "`threshold` means 10"},
		{"struct field", doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{1, 0}}, "has `income` as Currency"},
		{"inline rule", doc.StmtRef{Sec: codebook.SectionRules, Path: []int{0}}, "`r` when `f` equals 1"},
		{"case head", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}, "in case of `can go`"},
		{"require", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 0, 0}}, "`user` is required"},
		{"node", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 1, 0}}, "`ok` is a general rule"},
		{"attack", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 2, 0}}, "don't `ok` when `bad`"},
		{"exec", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0, 3, 0}}, "do `p`.`check` when `ok`"},
		{"endpoint head", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0}}, "provides `entry`"},
		{"endpoint require", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 0, 0}}, "`user` is required"},
		{"endpoint run", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 1, 0}}, "run `can go`"},
		{"endpoint check", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 2, 0}}, "check `can sub`"},
		{"internal head", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0}}, "on start"},
		{"internal run", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0, 1, 0}}, "run `can go`"},
		{"internal check", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0, 2, 0}}, "check `can sub`"},
	}
	for _, tc := range cases {
		got, err := (ENCodec{}).RenderStmt(d, tc.ref)
		if err != nil {
			t.Errorf("%s: RenderStmt: %v", tc.name, err)
			continue
		}
		if got != tc.want {
			t.Errorf("%s: RenderStmt = %q, want %q", tc.name, got, tc.want)
		}
	}
}
