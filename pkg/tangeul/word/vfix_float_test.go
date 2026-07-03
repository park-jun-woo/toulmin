//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestVFixFloat — VFix.Float returns the spec value for the golden bit patterns
package word

import "testing"

func TestVFixFloat(t *testing.T) {
	cases := []struct {
		in   VFix
		want float64
	}{
		{10000, 1.0},
		{-10000, -1.0},
		{8000, 0.8},
		{0, 0.0},
	}
	for _, c := range cases {
		if got := c.in.Float(); got != c.want {
			t.Fatalf("VFix(%d).Float()=%v want %v", c.in, got, c.want)
		}
	}
}
