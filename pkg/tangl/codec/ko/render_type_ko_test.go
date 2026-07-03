//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderTypeKo — renderTypeKo emits a bare built-in verbatim and a qualified ref as `a`.`N`
package ko

import "testing"

// TestRenderTypeKo checks that a built-in type passes through and a dotted ref
// becomes a backtick-qualified reference.
func TestRenderTypeKo(t *testing.T) {
	if got := renderTypeKo("Currency"); got != "Currency" {
		t.Errorf("built-in = %q", got)
	}
	if got := renderTypeKo("센서.Temperature"); got != "`센서`.`Temperature`" {
		t.Errorf("qualified = %q", got)
	}
}
