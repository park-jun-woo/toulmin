//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestStripLiteralQuotesKo — only one pair of matching outer double quotes is removed
package ko

import "testing"

func TestStripLiteralQuotesKo(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{`"값"`, "값"},
		{`""`, ""},
		{`"`, `"`},
		{`"값`, `"값`},
		{`""x""`, `"x"`},
		{"", ""},
	} {
		if got := stripLiteralQuotesKo(c.in); got != c.want {
			t.Errorf("stripLiteralQuotesKo(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
