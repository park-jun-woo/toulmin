//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestWordPrefix — a keyword must land on a word boundary to match
package en

import "testing"

// TestWordPrefix covers: no prefix match at all, an exact match (rest == ""), a
// match followed by a space, a match followed by a tab, and a false-positive
// guard where the text merely starts with the keyword as a substring of a longer
// word ("running" must not match "run").
func TestWordPrefix(t *testing.T) {
	cases := []struct {
		text, kw string
		want     bool
	}{
		{"check `x`", "run", false},
		{"run", "run", true},
		{"run `x`", "run", true},
		{"run\t`x`", "run", true},
		{"running `x`", "run", false},
	}
	for _, c := range cases {
		if got := wordPrefix(c.text, c.kw); got != c.want {
			t.Errorf("wordPrefix(%q, %q) = %v, want %v", c.text, c.kw, got, c.want)
		}
	}
}
