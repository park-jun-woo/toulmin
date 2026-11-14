//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSubject — renderSubject emits `this document is `x``
package en

import "testing"

// TestRenderSubject checks the subject declaration form.
func TestRenderSubject(t *testing.T) {
	if got := renderSubject("americano"); got != "this document is `americano`" {
		t.Errorf("renderSubject = %q", got)
	}
}
