//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidNumber — table-driven check of validNumber across float and non-float kinds
package types

import "testing"

func TestValidNumber(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"float32", float32(1.5), true},
		{"float64", float64(1.5), true},
		{"int", 3, false},
		{"string", "3.14", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validNumber(c.v)
			if got != c.want {
				t.Errorf("validNumber(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
