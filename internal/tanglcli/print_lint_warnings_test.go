//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestPrintLintWarnings — printLintWarnings writes each warning to w, one per line
package tanglcli

import (
	"bytes"
	"testing"
)

// TestPrintLintWarnings covers the single branch of printLintWarnings: the
// zero-iteration loop over an empty slice, and a multi-iteration loop over a
// non-empty slice.
func TestPrintLintWarnings(t *testing.T) {
	cases := []struct {
		name     string
		warnings []string
		want     string
	}{
		{
			name:     "no warnings writes nothing",
			warnings: nil,
			want:     "",
		},
		{
			name:     "warnings are written one per line",
			warnings: []string{"missing subject", "unknown alias"},
			want:     "missing subject\nunknown alias\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			printLintWarnings(&buf, c.warnings)
			if got := buf.String(); got != c.want {
				t.Errorf("printLintWarnings() output = %q, want %q", got, c.want)
			}
		})
	}
}
