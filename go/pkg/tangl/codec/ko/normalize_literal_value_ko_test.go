//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestNormalizeLiteralValueKo — quotes are stripped from the value, or from each item of an "is in" list
package ko

import "testing"

func TestNormalizeLiteralValueKo(t *testing.T) {
	for _, c := range []struct{ op, in, want string }{
		{"equals", `"VIP"`, "VIP"},
		{"equals", "plain", "plain"},
		{"is in", `"VIP", GOLD ,"x y"`, "VIP, GOLD, x y"},
		{"is in", "solo", "solo"},
	} {
		if got := normalizeLiteralValueKo(c.op, c.in); got != c.want {
			t.Errorf("normalizeLiteralValueKo(%q, %q) = %q, want %q", c.op, c.in, got, c.want)
		}
	}
}
