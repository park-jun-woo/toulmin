//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidDate — table-driven check of validDate across valid, malformed, and wrong-kind values
package types

import "testing"

func TestValidDate(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"valid date", "2026-07-01", true},
		{"bad format", "07/01/2026", false},
		{"wrong kind", 20260701, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validDate(c.v)
			if got != c.want {
				t.Errorf("validDate(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
