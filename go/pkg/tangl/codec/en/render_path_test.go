//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderPath — renderPath joins a single step and multi-step chains with dotted backticks
package en

import "testing"

// TestRenderPath checks that a nil step list renders empty (zero iterations),
// a single step renders as a lone backtick name (identical to a flat field),
// and multiple steps join with dots.
func TestRenderPath(t *testing.T) {
	cases := []struct {
		steps []string
		want  string
	}{
		{nil, ""},
		{[]string{"a"}, "`a`"},
		{[]string{"author", "trustScore"}, "`author`.`trustScore`"},
		{[]string{"a", "b", "c"}, "`a`.`b`.`c`"},
	}
	for _, tc := range cases {
		if got := renderPath(tc.steps); got != tc.want {
			t.Errorf("renderPath(%v) = %q, want %q", tc.steps, got, tc.want)
		}
	}
}
