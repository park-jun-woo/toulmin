//ff:func feature=tangl type=parser control=sequence
//ff:what TestSeeAnchorsKo — struck item skipped, non-struck indices count only parsed sees, parse error propagated
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestSeeAnchorsKo(t *testing.T) {
	t.Run("StruckSkippedThenTwoParsed", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "`old` 는 `old/pkg` 에서 참조한다", No: 1, Struck: true},
			{Text: "`a` 는 `a/pkg` 에서 참조한다", No: 2},
			{Text: "`b` 는 `b/pkg` 에서 참조한다", No: 3},
		}
		anchors, err := seeAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Sees) != 2 || a.Sees[0].Alias != "a" || a.Sees[1].Alias != "b" {
			t.Fatalf("sees = %+v", a.Sees)
		}
		if len(anchors) != 2 {
			t.Fatalf("anchors = %+v", anchors)
		}
		if anchors[0].Ref.Sec != codebook.SectionSee || len(anchors[0].Ref.Path) != 1 || anchors[0].Ref.Path[0] != 0 {
			t.Fatalf("anchors[0].Ref = %+v", anchors[0].Ref)
		}
		if anchors[1].Ref.Path[0] != 1 {
			t.Fatalf("anchors[1].Ref = %+v", anchors[1].Ref)
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "`a` 는 없음", No: 1}}
		if _, err := seeAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from parseSeeKo")
		}
	})
}
