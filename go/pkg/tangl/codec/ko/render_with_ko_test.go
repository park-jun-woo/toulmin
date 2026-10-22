//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestRenderWithKo — renderWithKo covers empty, single, multi 그리고 chains, and 받침 을/를 selection
package ko

import "testing"

// TestRenderWithKo checks the empty list (empty string), a single-term modifier,
// a multi-term 그리고 chain, and the 받침-driven object particle on the last term
// (록 받침 → 을, ASCII → 를, 력 받침 → 을).
func TestRenderWithKo(t *testing.T) {
	cases := []struct {
		with []string
		want string
	}{
		{nil, ""},
		{[]string{"차단 목록"}, "`차단 목록` 을 기준으로 하는"},
		{[]string{"a", "b", "c"}, "`a` 그리고 `b` 그리고 `c` 를 기준으로 하는"},
		{[]string{"신용 기준", "연체 이력"}, "`신용 기준` 그리고 `연체 이력` 을 기준으로 하는"},
	}
	for _, c := range cases {
		if got := renderWithKo(c.with); got != c.want {
			t.Errorf("renderWithKo(%v) = %q, want %q", c.with, got, c.want)
		}
	}
}
