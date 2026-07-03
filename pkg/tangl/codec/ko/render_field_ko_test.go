//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderFieldKo — renderFieldKo emits `` `f` 를 가진다, 타입 T `` with 받침-selected 을/를
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestRenderFieldKo checks the struct field form on a 받침 name (연소득→을) and an open
// name (신용점수→를), with a bare built-in type and a qualified ref type.
func TestRenderFieldKo(t *testing.T) {
	cases := []struct {
		f    ast.Field
		want string
	}{
		{ast.Field{Name: "연소득", Type: "Currency"}, "`연소득` 을 가진다, 타입 Currency"},
		{ast.Field{Name: "신용점수", Type: "Integer"}, "`신용점수` 를 가진다, 타입 Integer"},
		{ast.Field{Name: "등급", Type: "credit.Rating"}, "`등급` 을 가진다, 타입 `credit`.`Rating`"},
	}
	for _, c := range cases {
		if got := renderFieldKo(c.f); got != c.want {
			t.Errorf("renderFieldKo(%+v) = %q, want %q", c.f, got, c.want)
		}
	}
}
