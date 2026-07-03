//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestIndentWidth — tests indentWidth for spaces, tabs, mixed runs, non-whitespace, and empty input
package scan

import "testing"

func TestIndentWidth(t *testing.T) {
	cases := []struct {
		name string
		text string
		want int
	}{
		{"empty", "", 0},
		{"spaces only", "   ", 3},
		{"tab only", "\t", 4},
		{"mixed spaces and tabs", "  \t", 6},
		{"tab then default", "\ttext", 4},
		{"default immediately", "x", 0},
		{"space then default", " x", 1},
	}
	for _, c := range cases {
		if got := indentWidth(c.text); got != c.want {
			t.Errorf("%s: indentWidth(%q) = %d, want %d", c.name, c.text, got, c.want)
		}
	}
}
