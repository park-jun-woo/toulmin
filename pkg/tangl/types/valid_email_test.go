//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidEmail — table-driven check of validEmail across valid, malformed, and wrong-kind values
package types

import "testing"

func TestValidEmail(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"valid email", "user@example.com", true},
		{"missing at", "userexample.com", false},
		{"missing domain dot", "user@examplecom", false},
		{"wrong kind", 123, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validEmail(c.v)
			if got != c.want {
				t.Errorf("validEmail(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
