//ff:func feature=tangl type=parser control=sequence
//ff:what TestInternalAnchorsKo — struck skip, head/body success, and head/body error propagation
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestInternalAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "삭제됨", No: 1, Struck: true},
			{
				Text: "`주문 접수` 발생 시",
				No:   2,
				Span: scan.Span{Off: 0, Len: 20},
				Children: []scan.Item{
					{Text: "`물 붓기 가능` 진행한다", No: 3, Span: scan.Span{Off: 20, Len: 10}},
				},
			},
		}
		anchors, err := internalAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Internals) != 1 || a.Internals[0].Event != "주문 접수" {
			t.Fatalf("got internals %+v", a.Internals)
		}
		if len(anchors) != 2 {
			t.Fatalf("anchors = %d, want 2 (head + 1 body)", len(anchors))
		}
	})

	t.Run("HeadError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "정체불명의 문장", No: 1}}
		if _, err := internalAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error for a malformed internal head")
		}
	})

	t.Run("BodyError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{
				Text: "`이벤트` 발생 시",
				No:   1,
				Children: []scan.Item{
					{Text: "알 수 없는 문장", No: 2},
				},
			},
		}
		if _, err := internalAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from internalBodyAnchorsKo")
		}
	})
}
