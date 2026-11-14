//ff:func feature=tangl type=parser control=sequence
//ff:what TestCaseBodyAnchorsKo — struck skip, per-kind slot counters, and error propagation
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestCaseBodyAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c := &ast.Case{Name: "c"}
		children := []scan.Item{
			{Text: "삭제됨", No: 1, Struck: true},
			{Text: "`사용자` 는 필수다", No: 2, Span: scan.Span{Off: 0, Len: 10}},
			{Text: "`인증` 은 일반 규칙이다", No: 3, Span: scan.Span{Off: 10, Len: 10}},
			{Text: "`컵 없음` 이면 `주문 접수` 무효다", No: 4, Span: scan.Span{Off: 20, Len: 10}},
			{Text: "`주문 접수` 면 `물 붓기 가능` 진행한다", No: 5, Span: scan.Span{Off: 30, Len: 10}},
			{Text: "`컵 없음` 이면 `주문 접수` 배제한다", No: 6, Span: scan.Span{Off: 40, Len: 10}},
		}
		anchors, err := caseBodyAnchorsKo(c, 0, children, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(anchors) != 5 {
			t.Fatalf("anchors = %d, want 5", len(anchors))
		}
		if len(c.Requires) != 1 || len(c.Nodes) != 1 || len(c.Attacks) != 1 || len(c.Execs) != 1 || len(c.Supersedes) != 1 {
			t.Fatalf("counts req=%d node=%d atk=%d exec=%d sup=%d", len(c.Requires), len(c.Nodes), len(c.Attacks), len(c.Execs), len(c.Supersedes))
		}
		// slotKind order is Require·Node·Attack·Exec·Supersede (0..4), matching child order.
		for i, anchor := range anchors {
			if anchor.Ref.Path[0] != 0 || anchor.Ref.Path[1] != i || anchor.Ref.Path[2] != 0 {
				t.Errorf("anchor %d path = %v, want [0 %d 0]", i, anchor.Ref.Path, i)
			}
		}
	})

	t.Run("Error", func(t *testing.T) {
		c := &ast.Case{}
		children := []scan.Item{{Text: "알 수 없는 문장", No: 1}}
		if _, err := caseBodyAnchorsKo(c, 0, children, "t.md"); err == nil {
			t.Fatal("expected error for an unrecognized case-body statement")
		}
	})
}
