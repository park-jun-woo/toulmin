//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderSubjectKo — renderSubjectKo emits `이 문서는 `x` 이다`
package ko

import "testing"

// TestRenderSubjectKo checks the subject declaration form.
func TestRenderSubjectKo(t *testing.T) {
	if got := renderSubjectKo("아메리카노"); got != "이 문서는 `아메리카노` 이다" {
		t.Errorf("renderSubjectKo = %q", got)
	}
}
