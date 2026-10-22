//ff:func feature=policy type=engine control=iteration dimension=1
//ff:what TestFormatVerdict — verifies formatVerdict formats float64 to two decimal places
package policy

import "testing"

func TestFormatVerdict(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0.00"},
		{1, "1.00"},
		{-1, "-1.00"},
		{0.5555, "0.56"},
	}
	for _, c := range cases {
		if got := formatVerdict(c.in); got != c.want {
			t.Fatalf("formatVerdict(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
