//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCallPredKo — 참이면 predicate surface for alias/local FUNC_REF with 받침-selected 이/가
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestRenderCallPredKo(t *testing.T) {
	cases := []struct {
		c    ast.Call
		want string
	}{
		{ast.Call{Fn: ast.Ref{Alias: "검사", Name: "isSpam"}}, "`검사`.`isSpam` 가 참이면"},
		{ast.Call{Fn: ast.Ref{Name: "검증"}}, "`검증` 이 참이면"},
		{ast.Call{Fn: ast.Ref{Name: "점수화"}}, "`점수화` 가 참이면"},
	}
	for _, c := range cases {
		if got := renderCallPredKo(c.c); got != c.want {
			t.Errorf("renderCallPredKo(%#v) = %q, want %q", c.c, got, c.want)
		}
	}
}
