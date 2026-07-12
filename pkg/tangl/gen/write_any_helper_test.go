//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteAnyHelper — tests writeAnyHelper emits the tanglAny helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteAnyHelper(t *testing.T) {
	var w strings.Builder
	writeAnyHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglAny(src any, ok bool, where, expected string, elem func(any, int) (bool, any)) (bool, any) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "Actual: n}") {
		t.Errorf("expected examined-count Actual on the false verdict, got %q", got)
	}
}
