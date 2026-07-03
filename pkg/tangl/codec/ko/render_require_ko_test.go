//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderRequireKo — renderRequireKo covers the bare and 타입-annotated required-field forms
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderRequireKo checks the required-field form with 받침-selected 은/는 and an
// optional 타입 clause.
func TestRenderRequireKo(t *testing.T) {
	cases := []struct {
		r    ast.Require
		want string
	}{
		{ast.Require{Field: "사용자"}, "`사용자` 는 필수다"},
		{ast.Require{Field: "금액"}, "`금액` 은 필수다"},
		{ast.Require{Field: "사용자", Type: "Email"}, "`사용자` 는 필수다, 타입 Email"},
	}
	for _, c := range cases {
		if got := renderRequireKo(c.r); got != c.want {
			t.Errorf("renderRequireKo(%+v) = %q, want %q", c.r, got, c.want)
		}
	}
}
