//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseForceAttackKo — '… 를 F의 세기로 공격한다': success (attached/spaced 의, explicit 1.0), not-this-form, and error branches
package ko

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

func TestParseForceAttackKo(t *testing.T) {
	ok := []struct {
		text  string
		force float64
	}{
		{"`차단` 이면 `인증` 을 0.5의 세기로 공격한다", 0.5},
		{"`차단` 면 `인증` 를 0.25 의 세기로 공격한다", 0.25},
		{"`차단` 이면 `인증` 을 1.0의 세기로 공격한다", 1},
	}
	for _, tc := range ok {
		a, matched, err := parseForceAttackKo(scan.Item{Text: tc.text, No: 2}, "t.md")
		if err != nil || !matched || a.Attacker != "차단" || a.Target != "인증" || a.Force == nil || *a.Force != tc.force || a.Line != 2 {
			t.Errorf("parseForceAttackKo(%q) = (%+v, %v, %v)", tc.text, a, matched, err)
		}
	}
	if _, matched, err := parseForceAttackKo(scan.Item{Text: "`차단` 이면 `인증` 무효다", No: 1}, "t.md"); matched || err != nil {
		t.Errorf("무효다 must not match: matched %v err %v", matched, err)
	}
	errs := []struct {
		text, want string
	}{
		{"`차단` 이면 `인증` 을 0.5의 힘으로 공격한다", "expected '세기로'"},
		{"차단 이면 `인증` 을 0.5의 세기로 공격한다", "expected backtick-quoted attacker"},
		{"`차단` 이면 인증 을 0.5의 세기로 공격한다", "expected backtick-quoted target"},
		{"`차단` 이면 `인증` 을 1.5의 세기로 공격한다", "out of range"},
		{"`차단` 이면 `인증` 을 0.5의 강하게 세기로 공격한다", "unexpected trailing text"},
	}
	for _, tc := range errs {
		_, matched, err := parseForceAttackKo(scan.Item{Text: tc.text, No: 3}, "t.md")
		if !matched || err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("parseForceAttackKo(%q) = matched %v err %v, want %q", tc.text, matched, err, tc.want)
		}
	}
}
