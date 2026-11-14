//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderWith — renderWith covers the empty, single-term, and multi-term with/and chains
package en

import "testing"

// TestRenderWith checks the empty-list, single-term ("with" only), and
// multi-term ("with" then "and"...) forms of the with-clause loop.
func TestRenderWith(t *testing.T) {
	if got := renderWith(nil); got != "" {
		t.Errorf("empty: got %q, want %q", got, "")
	}
	if got := renderWith([]string{"a"}); got != " with `a`" {
		t.Errorf("single: got %q, want %q", got, " with `a`")
	}
	if got := renderWith([]string{"a", "b", "c"}); got != " with `a` and `b` and `c`" {
		t.Errorf("multi: got %q, want %q", got, " with `a` and `b` and `c`")
	}
}
