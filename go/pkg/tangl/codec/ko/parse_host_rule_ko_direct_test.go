//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseHostRuleKoDirect — a full host declaration yields name, host ref, criteria and grounds
package ko

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseHostRuleKoDirect(t *testing.T) {
	rule, ok, err := parseHostRuleKo(scan.Item{Text: "`시도 소진` 은 `최대 시도` 를 기준으로 하는 `시도 기록` 그리고 `로그` 를 근거로 하는 `래칫`.`triesExhausted` 쓰는 규칙이다", No: 7}, "t.md")
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	want := ast.InlineRule{Name: "시도 소진", Using: &ast.Ref{Alias: "래칫", Name: "triesExhausted"}, With: []string{"최대 시도"}, Grounds: []string{"시도 기록", "로그"}, Line: 7}
	if !reflect.DeepEqual(rule, want) {
		t.Fatalf("rule = %+v, want %+v", rule, want)
	}
	bare, ok, err := parseHostRuleKo(scan.Item{Text: "`단순` 은 `검사` 쓰는 규칙이다", No: 8}, "t.md")
	if err != nil || !ok || bare.Using == nil || bare.Using.Name != "검사" || bare.Using.Alias != "" || bare.Grounds != nil || bare.With != nil {
		t.Fatalf("bare = %+v ok=%v err=%v", bare, ok, err)
	}
}
