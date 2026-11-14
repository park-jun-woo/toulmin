//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderPathChainKo — backtick step joining over the 의 connective for empty/single/multi chains
package ko

import "testing"

func TestRenderPathChainKo(t *testing.T) {
	cases := []struct {
		steps []string
		want  string
	}{
		{nil, ""},
		{[]string{"기본가"}, "`기본가`"},
		{[]string{"작성자", "활성"}, "`작성자` 의 `활성`"},
		{[]string{"a", "b", "c"}, "`a` 의 `b` 의 `c`"},
	}
	for _, c := range cases {
		if got := renderPathChainKo(c.steps); got != c.want {
			t.Errorf("renderPathChainKo(%q) = %q, want %q", c.steps, got, c.want)
		}
	}
}
