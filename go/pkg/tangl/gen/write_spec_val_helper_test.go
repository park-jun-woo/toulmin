//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteSpecValHelper — tests writeSpecValHelper emits the tanglSpecVal helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteSpecValHelper(t *testing.T) {
	var w strings.Builder
	writeSpecValHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglSpecVal(specs toulmin.Specs, steps ...string) (any, bool) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "sp := specs.Find(steps[0])") {
		t.Errorf("expected Specs.Find lookup, got %q", got)
	}
}
