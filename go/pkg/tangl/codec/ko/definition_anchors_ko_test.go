//ff:func feature=tangl type=parser control=sequence
//ff:what TestDefinitionAnchorsKo — struck item skipped, const+struct defs anchor with field delegation, parse error propagated
package ko

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestDefinitionAnchorsKo(t *testing.T) {
	t.Run("StruckSkippedThenConstAndStruct", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "`old` 은 1 을 뜻한다", No: 1, Struck: true},
			{Text: "`문턱` 은 10 을 뜻한다", No: 2},
			{
				Text: "`신청인` 은 다음을 뜻한다",
				No:   3,
				Children: []scan.Item{
					{Text: "`연소득` 을 가진다, 타입 Currency", No: 4},
				},
			},
		}
		anchors, err := definitionAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Defs) != 2 || a.Defs[0].Name != "문턱" || a.Defs[1].Name != "신청인" {
			t.Fatalf("defs = %+v", a.Defs)
		}
		// def[0] head (di=0) + def[1] head (di=1) + def[1] field (di=1,fi=0)
		if len(anchors) != 3 {
			t.Fatalf("anchors = %+v", anchors)
		}
		if anchors[0].Ref.Sec != toulmin.SectionDefinitions || anchors[0].Ref.Path[0] != 0 {
			t.Fatalf("anchors[0].Ref = %+v", anchors[0].Ref)
		}
		if anchors[1].Ref.Path[0] != 1 {
			t.Fatalf("anchors[1].Ref = %+v", anchors[1].Ref)
		}
		if len(anchors[2].Ref.Path) != 2 || anchors[2].Ref.Path[0] != 1 || anchors[2].Ref.Path[1] != 0 {
			t.Fatalf("anchors[2].Ref = %+v", anchors[2].Ref)
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이름없음", No: 1}}
		if _, err := definitionAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from parseDefinitionKo")
		}
	})
}
