//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGroundsArg — tests groundsArg quoting for single, multiple, and empty name lists
package gen

import "testing"

func TestGroundsArg(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want string
	}{
		{"single", []string{"source-body"}, `"source-body"`},
		{"multiple preserves order", []string{"source-body", "mx"}, `"source-body", "mx"`},
		{"empty", nil, ""},
	}
	for _, c := range cases {
		if got := groundsArg(c.in); got != c.want {
			t.Errorf("groundsArg(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
