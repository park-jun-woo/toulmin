//ff:func feature=tangl type=parser control=sequence
//ff:what TestDefinitionFieldAnchorsKo — no children, struck field skipped, non-struck fields count with fi
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestDefinitionFieldAnchorsKo(t *testing.T) {
	t.Run("NoChildren", func(t *testing.T) {
		anchors := definitionFieldAnchorsKo(0, nil)
		if anchors != nil {
			t.Fatalf("anchors = %+v, want nil", anchors)
		}
	})

	t.Run("StruckSkippedThenTwoFields", func(t *testing.T) {
		children := []scan.Item{
			{Text: "`old` 을 가진다, 타입 Currency", No: 1, Struck: true},
			{Text: "`연소득` 을 가진다, 타입 Currency", No: 2},
			{Text: "`신용점수` 를 가진다, 타입 Integer", No: 3},
		}
		anchors := definitionFieldAnchorsKo(2, children)
		if len(anchors) != 2 {
			t.Fatalf("anchors = %+v", anchors)
		}
		if anchors[0].Ref.Sec != codebook.SectionDefinitions || len(anchors[0].Ref.Path) != 2 || anchors[0].Ref.Path[0] != 2 || anchors[0].Ref.Path[1] != 0 {
			t.Fatalf("anchors[0].Ref = %+v", anchors[0].Ref)
		}
		if anchors[1].Ref.Path[1] != 1 {
			t.Fatalf("anchors[1].Ref = %+v", anchors[1].Ref)
		}
	})
}
