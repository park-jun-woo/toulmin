//ff:what TestParseHostRuleKo — Korean Rules host declaration owns using and criteria
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestParseHostRuleKo(t *testing.T) {
	rule, err := parseInlineRuleKo(scan.Item{Text: "`시도 소진` 은 `최대 시도` 를 기준으로 하는 `시도 기록` 을 근거로 하는 `래칫`.`triesExhausted` 쓰는 규칙이다", No: 3}, "t.md")
	if err != nil {
		t.Fatal(err)
	}
	if rule.Using == nil || rule.Using.Alias != "래칫" || rule.Using.Name != "triesExhausted" || len(rule.With) != 1 || rule.With[0] != "최대 시도" || len(rule.Grounds) != 1 || rule.Grounds[0] != "시도 기록" {
		t.Fatalf("rule=%+v", rule)
	}
}
