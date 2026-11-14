//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestStripLiteralQuotes — only one pair of matching outer double quotes is removed
package parser

import "testing"

func TestStripLiteralQuotes(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{`"abc"`, "abc"},
		{`""`, ""},
		{`"`, `"`},
		{`"abc`, `"abc`},
		{`abc"`, `abc"`},
		{`""x""`, `"x"`},
		{"", ""},
	} {
		if got := stripLiteralQuotes(c.in); got != c.want {
			t.Errorf("stripLiteralQuotes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
