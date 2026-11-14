//ff:func feature=tangl type=parser control=sequence
//ff:what TestCaseAnchorsKo — struck skip, head/body success, and head/body error propagation
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestCaseAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "삭제된 사안", No: 1, Struck: true},
			{
				Text: "`컵 놓기 가능` 의 경우에",
				No:   2,
				Span: scan.Span{Off: 0, Len: 20},
				Children: []scan.Item{
					{Text: "`사용자` 는 필수다", No: 3, Span: scan.Span{Off: 20, Len: 10}},
				},
			},
		}
		anchors, err := caseAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Cases) != 1 || a.Cases[0].Name != "컵 놓기 가능" {
			t.Fatalf("got cases %+v", a.Cases)
		}
		if len(anchors) != 2 {
			t.Fatalf("anchors = %d, want 2 (head + 1 body)", len(anchors))
		}
		if anchors[0].Ref.Path[0] != 0 {
			t.Errorf("head anchor path = %v, want [0]", anchors[0].Ref.Path)
		}
	})

	t.Run("HeadError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "정체불명의 문장", No: 1}}
		if _, err := caseAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error for a malformed case head")
		}
	})

	t.Run("BodyError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{
				Text: "`X` 의 경우에",
				No:   1,
				Children: []scan.Item{
					{Text: "알 수 없는 문장", No: 2},
				},
			},
		}
		if _, err := caseAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from caseBodyAnchorsKo")
		}
	})
}
