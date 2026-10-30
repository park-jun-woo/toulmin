//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseHostRuleKoErrors — non-declarations are unmatched; malformed names, modifiers, refs, checking, weight and level fail
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseHostRuleKoErrors(t *testing.T) {
	if _, ok, err := parseHostRuleKo(scan.Item{Text: "`점수` 가 5 보다 크면 `합격` 이다", No: 1}, "t.md"); ok || err != nil {
		t.Fatalf("condition rule: ok=%v err=%v, want unmatched", ok, err)
	}
	for _, c := range []struct{ text, want string }{
		{"이름 없는 규칙이다", "expected backtick-quoted rule name"},
		{"`r` 은 가중치 abc 의 `래칫`.`f` 쓰는 규칙이다", "expected float after '가중치'"},
		{"`r` 은 `래칫`.`f` 규칙이다", "expected function reference followed by '쓰는'"},
		{"`r` 은 그냥 쓰는 규칙이다", "expected function reference before '쓰는'"},
		{"`r` 은 `사안` 판정 쓰는 규칙이다", "requires a host function, not a checking case"},
		{"`r` 은 가중치 0.5 의 `래칫`.`f` 쓰는 규칙이다", "weight and level belong to a case binding"},
		{"`r` 은 실패를 보고하는 `래칫`.`f` 쓰는 규칙이다", "weight and level belong to a case binding"},
	} {
		_, ok, err := parseHostRuleKo(scan.Item{Text: c.text, No: 2}, "t.md")
		if !ok || err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%q: ok=%v err=%v, want %q", c.text, ok, err, c.want)
		}
	}
}
