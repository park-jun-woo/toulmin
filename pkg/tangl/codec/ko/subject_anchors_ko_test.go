//ff:func feature=tangl type=parser control=sequence
//ff:what TestSubjectAnchorsKo — struck item skipped, valid item sets Subject and anchors, parse error propagated
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestSubjectAnchorsKo(t *testing.T) {
	t.Run("StruckSkipped", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이 문서는 `x` 이다", No: 1, Struck: true}}
		anchors, err := subjectAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if anchors != nil || a.Subject != "" {
			t.Fatalf("anchors=%+v subject=%q, want empty", anchors, a.Subject)
		}
	})

	t.Run("Valid", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이 문서는 `아메리카노` 이다", No: 1}}
		anchors, err := subjectAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.Subject != "아메리카노" {
			t.Fatalf("subject = %q", a.Subject)
		}
		if len(anchors) != 1 || anchors[0].Ref.Sec != codebook.SectionSubject || anchors[0].Ref.Path != nil {
			t.Fatalf("anchors = %+v", anchors)
		}
	})

	t.Run("ParseError", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이 문서는 x 이다", No: 1}}
		if _, err := subjectAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error propagated from parseSubjectKo")
		}
	})
}
