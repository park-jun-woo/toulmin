//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteCtxValHelper — tests writeCtxValHelper emits the tanglCtxVal helper source
package gen

import (
	"strings"
	"testing"
)

func TestWriteCtxValHelper(t *testing.T) {
	var w strings.Builder
	writeCtxValHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglCtxVal(ctx toulmin.Context, steps ...string) (any, bool) {") {
		t.Errorf("expected helper function signature, got %q", got)
	}
	if !strings.Contains(got, "return tanglWalk(v, steps[1:]...)") {
		t.Errorf("expected remaining-steps walk, got %q", got)
	}
}
