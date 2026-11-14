//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidCurrency — table-driven check of validCurrency across numeric and non-numeric kinds
package types

import "testing"

func TestValidCurrency(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"int", 1000, true},
		{"float", 19.99, true},
		{"string", "19.99", false},
		{"bool", true, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validCurrency(c.v)
			if got != c.want {
				t.Errorf("validCurrency(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
