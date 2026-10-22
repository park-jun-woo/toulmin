//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestFormatFloat — tests formatFloat renders the shortest round-tripping literal
package gen

import "testing"

func TestFormatFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{1, "1"},
		{0.5, "0.5"},
		{-1.25, "-1.25"},
	}
	for _, c := range cases {
		if got := formatFloat(c.in); got != c.want {
			t.Errorf("formatFloat(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
