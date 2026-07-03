//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderNodeKo — renderNodeKo covers role words, using/checking, and 기준/가중 tails
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderNodeKo checks the node variants: a bare general rule, a counter rule
// using a function with a 기준 clause, an except rule with 기준+가중, and a general rule
// composing another case via 판정 (checking).
func TestRenderNodeKo(t *testing.T) {
	q := 0.8
	cases := []struct {
		n    ast.Node
		want string
	}{
		{ast.Node{Name: "인증", Role: ast.GeneralRule},
			"`인증` 은 일반 규칙이다"},
		{ast.Node{Name: "IP 차단", Role: ast.CounterRule, Using: &ast.Ref{Alias: "정책", Name: "IsIPBlocked"}, With: []string{"차단 목록"}},
			"`IP 차단` 은 `정책`.`IsIPBlocked` 쓰는 반박 규칙이다, 기준 `차단 목록`"},
		{ast.Node{Name: "저온", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "센서", Name: "low"}, With: []string{"a", "b"}, Qualified: &q},
			"`저온` 은 `센서`.`low` 쓰는 예외 규칙이다, 기준 `a` 그리고 `b`, 가중 0.8"},
		{ast.Node{Name: "승인", Role: ast.GeneralRule, Checking: "승인 가능"},
			"`승인` 은 `승인 가능` 판정 쓰는 일반 규칙이다"},
		{ast.Node{Name: "컵 없음", Role: ast.CounterRule, Using: &ast.Ref{Alias: "게이트", Name: "noCup"}, Level: ast.LevelFail},
			"`컵 없음` 은 `게이트`.`noCup` 쓰는 반박 규칙이다, 수준 실패"},
		{ast.Node{Name: "재고 부족", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "게이트", Name: "lowStock"}, Level: ast.LevelReview},
			"`재고 부족` 은 `게이트`.`lowStock` 쓰는 예외 규칙이다, 수준 검토"},
	}
	for _, c := range cases {
		if got := renderNodeKo(c.n); got != c.want {
			t.Errorf("renderNodeKo(%+v) = %q, want %q", c.n, got, c.want)
		}
	}
}
