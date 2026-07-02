//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteCompareHelper — tests writeCompareHelper emits the tanglCompare helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteCompareHelper(t *testing.T) {
	var w strings.Builder
	writeCompareHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglCompare(ctx toulmin.Context, field, op, value string) bool {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, `case "is empty":`) {
		t.Errorf("expected is empty case, got %q", got)
	}
	if !strings.Contains(got, `case "is greater than", "is less than", "is at most", "is at least":`) {
		t.Errorf("expected numeric comparison case, got %q", got)
	}
}
