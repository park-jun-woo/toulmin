//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidQuantity — table-driven check of validQuantity across numeric and non-numeric kinds
package types

import "testing"

func TestValidQuantity(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"int", 5, true},
		{"float", 5.5, true},
		{"string", "5", false},
		{"bool", false, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validQuantity(c.v)
			if got != c.want {
				t.Errorf("validQuantity(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
