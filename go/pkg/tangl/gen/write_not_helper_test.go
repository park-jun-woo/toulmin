//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteNotHelper — tests writeNotHelper emits the tanglNot helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteNotHelper(t *testing.T) {
	var w strings.Builder
	writeNotHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglNot(term func() (bool, any)) (bool, any) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "return !b, w") {
		t.Errorf("expected verdict inversion, got %q", got)
	}
}
