//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderCertaintyKo — renderCertaintyKo maps each canonical threshold Op to its 확신 token
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestRenderCertaintyKo checks the four confidence-gate thresholds and the
// unknown-operator error.
func TestRenderCertaintyKo(t *testing.T) {
	cases := []struct {
		c    ast.Certainty
		want string
	}{
		{ast.Certainty{Op: "at least", Percent: 75}, "확신 75% 이상이면"},
		{ast.Certainty{Op: "above", Percent: 50}, "확신 50% 초과면"},
		{ast.Certainty{Op: "less than", Percent: 20}, "확신 20% 미만이면"},
		{ast.Certainty{Op: "at most", Percent: 90}, "확신 90% 이하면"},
	}
	for _, c := range cases {
		got, err := renderCertaintyKo(c.c)
		if err != nil {
			t.Errorf("renderCertaintyKo(%+v): %v", c.c, err)
			continue
		}
		if got != c.want {
			t.Errorf("renderCertaintyKo(%+v) = %q, want %q", c.c, got, c.want)
		}
	}
	if _, err := renderCertaintyKo(ast.Certainty{Op: "bogus", Percent: 1}); err == nil {
		t.Error("renderCertaintyKo accepted an unknown operator")
	}
}
