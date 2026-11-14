//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderExecKo — renderExecKo covers do/한 번만/확신, undo 로/으로, and run edges
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderExecKo checks the do (plain, 한 번만, 한 번만+확신), undo (with 받침-selected
// 로/으로), and run forms, plus the unknown-kind error.
func TestRenderExecKo(t *testing.T) {
	arm := &ast.Ref{Alias: "팔", Name: "placeCup"}
	cases := []struct {
		e    ast.Exec
		want string
	}{
		{ast.Exec{Kind: ast.DoExec, Func: arm, Node: "주문 접수"},
			"`주문 접수` 면 `팔`.`placeCup` 실행한다"},
		{ast.Exec{Kind: ast.DoExec, Func: arm, Node: "주문 접수", Once: true},
			"`주문 접수` 면 `팔`.`placeCup` 한 번만 실행한다"},
		{ast.Exec{Kind: ast.DoExec, Func: arm, Node: "주문 접수", Once: true, Certainty: &ast.Certainty{Op: "at least", Percent: 75}},
			"`주문 접수` 확신 75% 이상이면 `팔`.`placeCup` 한 번만 실행한다"},
		{ast.Exec{Kind: ast.UndoExec, Func: &ast.Ref{Alias: "은행", Name: "refund"}, Node: "잔액 충분"},
			"`잔액 충분` 이면 `은행`.`refund` 로 되돌린다"},
		{ast.Exec{Kind: ast.RunExec, Case: "입금 가능", Node: "잔액 충분"},
			"`잔액 충분` 이면 `입금 가능` 진행한다"},
	}
	for _, c := range cases {
		got, err := renderExecKo(c.e)
		if err != nil {
			t.Errorf("renderExecKo(%+v): %v", c.e, err)
			continue
		}
		if got != c.want {
			t.Errorf("renderExecKo(%+v) = %q, want %q", c.e, got, c.want)
		}
	}
	if _, err := renderExecKo(ast.Exec{Kind: ast.ExecKind(99), Node: "n"}); err == nil {
		t.Error("renderExecKo accepted an unknown ExecKind")
	}
	if _, err := renderExecKo(ast.Exec{Kind: ast.DoExec, Func: arm, Node: "주문 접수", Certainty: &ast.Certainty{Op: "bogus", Percent: 1}}); err == nil {
		t.Error("renderExecKo accepted a DoExec with an unknown certainty operator")
	}
}
