//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestValidURL — table-driven check of validURL across valid, malformed, and wrong-kind values
package types

import "testing"

func TestValidURL(t *testing.T) {
	cases := []struct {
		name string
		v    any
		want bool
	}{
		{"valid url", "https://example.com/path", true},
		{"parse error", "example.com/path", false},
		{"missing host", "mailto:foo@bar.com", false},
		{"missing scheme and host", "/relative/path", false},
		{"wrong kind", 123, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := validURL(c.v)
			if got != c.want {
				t.Errorf("validURL(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}
