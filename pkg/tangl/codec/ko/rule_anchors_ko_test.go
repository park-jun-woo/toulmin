//ff:func feature=tangl type=parser control=sequence
//ff:what TestRuleAnchorsKo — struck skip, whole-subtree span anchoring, and parse-error propagation
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestRuleAnchorsKo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{
			{Text: "삭제됨", No: 1, Struck: true},
			{
				Text: "`x` `n1` 이다",
				No:   2,
				Span: scan.Span{Off: 0, Len: 20},
			},
		}
		anchors, err := ruleAnchorsKo(a, items, "t.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(a.Rules) != 1 || a.Rules[0].Name != "n1" {
			t.Fatalf("got rules %+v", a.Rules)
		}
		if len(anchors) != 1 || anchors[0].Span.Off != 0 || anchors[0].Span.Len != 20 {
			t.Fatalf("anchors = %+v", anchors)
		}
	})

	t.Run("Error", func(t *testing.T) {
		a := &ast.Document{}
		items := []scan.Item{{Text: "이상한 문장", No: 1}}
		if _, err := ruleAnchorsKo(a, items, "t.md"); err == nil {
			t.Fatal("expected error for a malformed rule")
		}
	})
}
