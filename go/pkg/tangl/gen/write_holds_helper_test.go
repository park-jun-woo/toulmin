//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteHoldsHelper — tests writeHoldsHelper emits the tanglHolds helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteHoldsHelper(t *testing.T) {
	var w strings.Builder
	writeHoldsHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglHolds(b bool, ev any, where string) (bool, any) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, `Expected: "holds"`) {
		t.Errorf("expected holds Fact fallback, got %q", got)
	}
}
