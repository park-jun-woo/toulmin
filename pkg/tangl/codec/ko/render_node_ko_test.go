//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderNodeKo — renderNodeKo covers role words, using/checking, and 0~3 adnominal modifiers
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderNodeKo checks the node variants in the adnominal canonical form: a bare
// general rule, a single-기준 counter rule using a function, a 기준+가중 except rule, a
// checking (판정) general rule, a 보고 (fail/review) level rule, the three-modifier
// composition (기준, 가중, 보고 joined by ", "), and a bare-기준 general rule.
func TestRenderNodeKo(t *testing.T) {
	q08 := 0.8
	cases := []struct {
		n    ast.Node
		want string
	}{
		{ast.Node{Name: "인증", Role: ast.GeneralRule},
			"`인증` 은 일반 규칙이다"},
		{ast.Node{Name: "IP 차단", Role: ast.CounterRule, Using: &ast.Ref{Alias: "정책", Name: "IsIPBlocked"}, With: []string{"차단 목록"}},
			"`IP 차단` 은 `차단 목록` 을 기준으로 하는 `정책`.`IsIPBlocked` 쓰는 반박 규칙이다"},
		{ast.Node{Name: "저온", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "센서", Name: "low"}, With: []string{"a", "b"}, Qualified: &q08},
			"`저온` 은 `a` 그리고 `b` 를 기준으로 하는, 가중치 0.8 의 `센서`.`low` 쓰는 예외 규칙이다"},
		{ast.Node{Name: "승인", Role: ast.GeneralRule, Checking: "승인 가능"},
			"`승인` 은 `승인 가능` 판정 쓰는 일반 규칙이다"},
		{ast.Node{Name: "컵 없음", Role: ast.CounterRule, Using: &ast.Ref{Alias: "게이트", Name: "noCup"}, Level: ast.LevelFail},
			"`컵 없음` 은 실패를 보고하는 `게이트`.`noCup` 쓰는 반박 규칙이다"},
		{ast.Node{Name: "재고 부족", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "게이트", Name: "lowStock"}, Level: ast.LevelReview},
			"`재고 부족` 은 검토를 보고하는 `게이트`.`lowStock` 쓰는 예외 규칙이다"},
		{ast.Node{Name: "IP 차단", Role: ast.CounterRule, Using: &ast.Ref{Alias: "정책", Name: "IsIPBlocked"}, With: []string{"차단 목록"}, Qualified: &q08, Level: ast.LevelFail},
			"`IP 차단` 은 `차단 목록` 을 기준으로 하는, 가중치 0.8 의, 실패를 보고하는 `정책`.`IsIPBlocked` 쓰는 반박 규칙이다"},
		{ast.Node{Name: "기본 승인", Role: ast.GeneralRule, With: []string{"승인 기준"}},
			"`기본 승인` 은 `승인 기준` 을 기준으로 하는 일반 규칙이다"},
	}
	for _, c := range cases {
		if got := renderNodeKo(c.n); got != c.want {
			t.Errorf("renderNodeKo(%+v) = %q, want %q", c.n, got, c.want)
		}
	}
}
