//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestTrimBackticks — trimBackticks strips leading and trailing backticks
package tanglcli

import "testing"

// TestTrimBackticks covers trimBackticks's single unconditional branch
// across inputs with and without surrounding backticks.
func TestTrimBackticks(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "no backticks", in: "transfer", want: "transfer"},
		{name: "surrounded by backticks", in: "`transfer`", want: "transfer"},
		{name: "empty string", in: "", want: ""},
		{name: "only backticks", in: "``", want: ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := trimBackticks(c.in); got != c.want {
				t.Errorf("trimBackticks(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
