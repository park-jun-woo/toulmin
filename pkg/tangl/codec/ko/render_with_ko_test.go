//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderWithKo — renderWithKo covers empty, single, and multi-item 기준/그리고 tails
package ko

import "testing"

// TestRenderWithKo checks the empty tail (no loop iterations), a single-item 기준
// clause (the i==0 branch only), and a multi-item clause chaining 기준 then 그리고 for
// every following item (the i==0 and i!=0 branches both firing).
func TestRenderWithKo(t *testing.T) {
	cases := []struct {
		with []string
		want string
	}{
		{nil, ""},
		{[]string{"차단 목록"}, ", 기준 `차단 목록`"},
		{[]string{"a", "b", "c"}, ", 기준 `a` 그리고 `b` 그리고 `c`"},
	}
	for _, c := range cases {
		if got := renderWithKo(c.with); got != c.want {
			t.Errorf("renderWithKo(%v) = %q, want %q", c.with, got, c.want)
		}
	}
}
