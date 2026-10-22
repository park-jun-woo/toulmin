//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidTime — table-driven check of validTime across valid, malformed, and wrong-kind values
package types

import "testing"

func TestValidTime(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"hms ok", "13:45:00", true},
		{"hm ok", "13:45", true},
		{"bad format", "1pm", false},
		{"wrong kind", 1345, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validTime(c.v)
			if got != c.want {
				t.Errorf("validTime(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
