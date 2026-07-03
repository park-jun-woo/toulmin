//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseRoleKo — 일반/반박/예외 map to their roles, unknown word rejected
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseRoleKo(t *testing.T) {
	cases := []struct {
		word string
		want ast.Role
	}{
		{"일반", ast.GeneralRule},
		{"반박", ast.CounterRule},
		{"예외", ast.ExceptRule},
	}
	for _, c := range cases {
		got, ok := parseRoleKo(c.word)
		if !ok || got != c.want {
			t.Errorf("parseRoleKo(%q) = %v ok=%v, want %v", c.word, got, ok, c.want)
		}
	}
	if _, ok := parseRoleKo("특수"); ok {
		t.Error("expected ok=false for an unknown role word")
	}
}
