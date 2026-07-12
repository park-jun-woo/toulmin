//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteAndHelper — tests writeAndHelper emits the tanglAnd helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteAndHelper(t *testing.T) {
	var w strings.Builder
	writeAndHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglAnd(terms ...func() (bool, any)) (bool, any) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "return false, ww") {
		t.Errorf("expected first-false short circuit, got %q", got)
	}
}
