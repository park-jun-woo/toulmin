//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidBoolean — table-driven check of validBoolean across bool and non-bool kinds
package types

import "testing"

func TestValidBoolean(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"string", "true", false},
		{"int", 1, false},
		{"nil", nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validBoolean(c.v)
			if got != c.want {
				t.Errorf("validBoolean(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
