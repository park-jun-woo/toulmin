//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRunCheckKo — renderRunCheckKo emits 진행한다 (run) / 판정한다 (check) by slotKind
package ko

import "testing"

// TestRenderRunCheckKo checks the run (slotKind 1) and check (slotKind 2) body forms
// and the unknown-slotKind error.
func TestRenderRunCheckKo(t *testing.T) {
	if got, err := renderRunCheckKo(1, "컵 놓기 가능"); err != nil || got != "`컵 놓기 가능` 진행한다" {
		t.Errorf("run = %q, err %v", got, err)
	}
	if got, err := renderRunCheckKo(2, "접근 가능"); err != nil || got != "`접근 가능` 판정한다" {
		t.Errorf("check = %q, err %v", got, err)
	}
	if _, err := renderRunCheckKo(0, "x"); err == nil {
		t.Error("renderRunCheckKo accepted an unknown slotKind")
	}
}
