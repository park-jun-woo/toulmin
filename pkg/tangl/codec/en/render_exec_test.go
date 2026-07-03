//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderExec — renderExec covers do/once/certainty, undo, and run edges
package en

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderExec checks the do (plain, once, once+certainty), undo, and run forms.
func TestRenderExec(t *testing.T) {
	fn := &ast.Ref{Alias: "arm", Name: "placeCup"}
	cases := []struct {
		e    ast.Exec
		want string
	}{
		{ast.Exec{Kind: ast.DoExec, Func: fn, Node: "n"}, "do `arm`.`placeCup` when `n`"},
		{ast.Exec{Kind: ast.DoExec, Func: fn, Node: "n", Once: true}, "do `arm`.`placeCup` once when `n`"},
		{ast.Exec{Kind: ast.DoExec, Func: fn, Node: "n", Once: true, Certainty: &ast.Certainty{Op: "at least", Percent: 75}},
			"do `arm`.`placeCup` once when `n` if at least 75% certain"},
		{ast.Exec{Kind: ast.UndoExec, Func: &ast.Ref{Alias: "bank", Name: "refund"}, Node: "ok"},
			"undo `bank`.`refund` when `ok`"},
		{ast.Exec{Kind: ast.RunExec, Case: "can deposit", Node: "ok"},
			"run `can deposit` when `ok`"},
	}
	for _, tc := range cases {
		got, err := renderExec(tc.e)
		if err != nil {
			t.Errorf("renderExec(%+v): %v", tc.e, err)
			continue
		}
		if got != tc.want {
			t.Errorf("renderExec = %q, want %q", got, tc.want)
		}
	}

	if _, err := renderExec(ast.Exec{Kind: ast.ExecKind(99), Node: "n"}); err == nil {
		t.Error("renderExec accepted an unknown ExecKind")
	}
}
