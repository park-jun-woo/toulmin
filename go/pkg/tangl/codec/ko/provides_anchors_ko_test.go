//ff:func feature=tangl type=parser control=sequence
//ff:what TestProvidesAnchorsKo — struck item skipped, valid endpoint with body delegation, head parse error, body error propagated
package ko

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestProvidesAnchorsKo(t *testing.T) {
	t.Run("StruckSkippedThenValidWithBody", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "`old` 를 한다면", No: 1, Struck: true},
			{
				Text: "`주문 접수` 를 한다면",
				No:   2,
				Children: []scan.Item{
					{Text: "`사용자` 는 필수다", No: 3},
					{Text: "`컵 놓기 가능` 진행한다", No: 4},
				},
			},
		}
		anchors, err := providesAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Provides) != 1 || a.Provides[0].Name != "주문 접수" {
			t.Fatalf("provides = %+v", a.Provides)
		}
		if len(a.Provides[0].Requires) != 1 || len(a.Provides[0].Runs) != 1 {
			t.Fatalf("endpoint body = %+v", a.Provides[0])
		}
		// head anchor (ei=0) + require body anchor + run body anchor
		if len(anchors) != 3 {
			t.Fatalf("anchors = %+v", anchors)
		}
		if anchors[0].Ref.Sec != toulmin.SectionProvides || len(anchors[0].Ref.Path) != 1 || anchors[0].Ref.Path[0] != 0 {
			t.Fatalf("anchors[0].Ref = %+v", anchors[0].Ref)
		}
		if len(anchors[1].Ref.Path) != 3 || anchors[1].Ref.Path[1] != 0 {
			t.Fatalf("anchors[1].Ref = %+v", anchors[1].Ref)
		}
		if len(anchors[2].Ref.Path) != 3 || anchors[2].Ref.Path[1] != 1 {
			t.Fatalf("anchors[2].Ref = %+v", anchors[2].Ref)
		}
	})

	t.Run("HeadParseError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이름없음", No: 1}}
		if _, err := providesAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from parseEndpointKo")
		}
	})

	t.Run("BodyError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{
				Text: "`주문 접수` 를 한다면",
				No:   1,
				Children: []scan.Item{
					{Text: "알 수 없는 문장", No: 2},
				},
			},
		}
		if _, err := providesAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from endpointBodyAnchorsKo")
		}
	})
}
