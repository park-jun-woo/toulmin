//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderInlineRuleKoHost — a host rule renders '… 쓰는 규칙이다' with optional 기준/근거 modifiers and re-parses
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestRenderInlineRuleKoHost(t *testing.T) {
	bare := ast.InlineRule{Name: "시도 소진", Using: &ast.Ref{Alias: "래칫", Name: "triesExhausted"}}
	got, err := renderInlineRuleKo(bare)
	if err != nil || got != "`시도 소진` 은 `래칫`.`triesExhausted` 쓰는 규칙이다" {
		t.Fatalf("bare = %q, %v", got, err)
	}
	full := ast.InlineRule{Name: "시도 소진", Using: &ast.Ref{Alias: "래칫", Name: "triesExhausted"}, With: []string{"최대 시도"}, Grounds: []string{"시도 기록"}, Line: 3}
	got, err = renderInlineRuleKo(full)
	if err != nil || got != "`시도 소진` 은 `최대 시도` 를 기준으로 하는 `시도 기록` 을 근거로 하는 `래칫`.`triesExhausted` 쓰는 규칙이다" {
		t.Fatalf("full = %q, %v", got, err)
	}
	back, err := parseInlineRuleKo(scan.Item{Text: got, No: 3}, "t.md")
	if err != nil || !reflect.DeepEqual(back, full) {
		t.Fatalf("round trip = %+v, %v", back, err)
	}
}
