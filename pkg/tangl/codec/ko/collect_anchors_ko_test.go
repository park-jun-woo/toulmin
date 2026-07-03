//ff:func feature=tangl type=parser control=sequence
//ff:what TestCollectAnchorsKo — empty input, unknown heading, out-of-order sections, section-body error, valid multi-section walk
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestCollectAnchorsKo(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		a := &ast.Document{}
		anchors, err := collectAnchorsKo(a, nil, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if anchors != nil {
			t.Fatalf("anchors = %+v, want nil", anchors)
		}
	})

	t.Run("UnknownSection", func(t *testing.T) {
		secs := []scan.Section{{Name: "모름"}}
		a := &ast.Document{}
		if _, err := collectAnchorsKo(a, secs, "t.md"); err == nil || !strings.Contains(err.Error(), "unknown tangl section") {
			t.Fatalf("expected unknown-section error, got %v", err)
		}
	})

	t.Run("OutOfOrder", func(t *testing.T) {
		secs := []scan.Section{
			{Name: "참조", Lines: []scan.Line{{Text: "- `a` 는 `b` 에서 참조한다", No: 1}}},
			{Name: "주제", Lines: []scan.Line{{Text: "- 이 문서는 `x` 이다", No: 2}}},
		}
		a := &ast.Document{}
		if _, err := collectAnchorsKo(a, secs, "t.md"); err == nil || !strings.Contains(err.Error(), "out of order") {
			t.Fatalf("expected out-of-order error, got %v", err)
		}
	})

	t.Run("SectionBodyError", func(t *testing.T) {
		secs := []scan.Section{
			{Name: "주제", Lines: []scan.Line{{Text: "- 이 문서는 x 이다", No: 1}}},
		}
		a := &ast.Document{}
		if _, err := collectAnchorsKo(a, secs, "t.md"); err == nil {
			t.Fatal("expected error propagated from sectionAnchorsKo")
		}
	})

	t.Run("ValidMultiSection", func(t *testing.T) {
		secs := []scan.Section{
			{Name: "주제", Lines: []scan.Line{{Text: "- 이 문서는 `아메리카노` 이다", No: 1}}},
			{Name: "참조", Lines: []scan.Line{{Text: "- `a` 는 `b` 에서 참조한다", No: 2}}},
		}
		a := &ast.Document{}
		anchors, err := collectAnchorsKo(a, secs, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.Subject != "아메리카노" || len(a.Sees) != 1 {
			t.Fatalf("got subject=%q sees=%+v", a.Subject, a.Sees)
		}
		// 2 SectionBlocks + 2 StmtBlocks
		if len(anchors) != 4 {
			t.Fatalf("anchors = %+v, want 4", anchors)
		}
	})
}
