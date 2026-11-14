//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestNormalizeLiteralValue — quotes are stripped from the value, or from each item of an "is in" list
package parser

import "testing"

func TestNormalizeLiteralValue(t *testing.T) {
	for _, c := range []struct{ op, in, want string }{
		{"equals", `"v0.3.1"`, "v0.3.1"},
		{"equals", "plain", "plain"},
		{"equals", `"a", "b"`, `a", "b`},
		{"is in", `"VIP", GOLD ,"x y"`, "VIP, GOLD, x y"},
		{"is in", "solo", "solo"},
	} {
		if got := normalizeLiteralValue(c.op, c.in); got != c.want {
			t.Errorf("normalizeLiteralValue(%q, %q) = %q, want %q", c.op, c.in, got, c.want)
		}
	}
}
