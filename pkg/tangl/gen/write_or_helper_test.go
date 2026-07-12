//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteOrHelper — tests writeOrHelper emits the tanglOr helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteOrHelper(t *testing.T) {
	var w strings.Builder
	writeOrHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglOr(terms ...func() (bool, any)) (bool, any) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "return true, ww") {
		t.Errorf("expected first-true short circuit, got %q", got)
	}
}
