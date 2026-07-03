//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderType — renderType emits bare built-ins and qualified backtick refs
package en

import "testing"

// TestRenderType checks a bare built-in token and a qualified reference token.
func TestRenderType(t *testing.T) {
	if got := renderType("Currency"); got != "Currency" {
		t.Errorf("builtin = %q", got)
	}
	if got := renderType("credit.Threshold"); got != "`credit`.`Threshold`" {
		t.Errorf("qualified = %q", got)
	}
}
