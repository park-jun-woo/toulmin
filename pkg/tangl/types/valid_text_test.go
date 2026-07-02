//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidText — table-driven check of validText across string and non-string kinds
package types

import "testing"

func TestValidText(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"string", "hello", true},
		{"empty string", "", true},
		{"int", 42, false},
		{"nil", nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validText(c.v)
			if got != c.want {
				t.Errorf("validText(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
