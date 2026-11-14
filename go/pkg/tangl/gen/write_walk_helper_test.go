//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteWalkHelper — tests writeWalkHelper emits the tanglWalk and tanglField sources
package gen

import (
	"strings"
	"testing"
)

func TestWriteWalkHelper(t *testing.T) {
	var w strings.Builder
	writeWalkHelper(&w)
	got := w.String()
	if !strings.Contains(got, "func tanglWalk(cur any, steps ...string) (any, bool) {") {
		t.Errorf("expected tanglWalk signature, got %q", got)
	}
	if !strings.Contains(got, "func tanglField(v reflect.Value, name string) (reflect.Value, bool) {") {
		t.Errorf("expected tanglField signature, got %q", got)
	}
	if !strings.Contains(got, "strings.EqualFold(sf.Name, name)") {
		t.Errorf("expected case-insensitive field match, got %q", got)
	}
}
