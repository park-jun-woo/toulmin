//ff:func feature=engine type=engine control=sequence
//ff:what TestSafeCallHandler — safeCallHandler returns nil, the handler error, or a recovered panic
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

func TestSafeCallHandler(t *testing.T) {
	self := TraceEntry{Name: "n"}
	tr := Trace{nodes: []TraceEntry{self}, ctx: NewContext()}

	if err := safeCallHandler(func(self TraceEntry, t Trace) error { return nil }, self, tr); err != nil {
		t.Errorf("ok handler should return nil, got %v", err)
	}

	want := fmt.Errorf("bad")
	if err := safeCallHandler(func(self TraceEntry, t Trace) error { return want }, self, tr); err != want {
		t.Errorf("error handler should propagate error, got %v", err)
	}

	err := safeCallHandler(func(self TraceEntry, t Trace) error { panic("boom") }, self, tr)
	if err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Errorf("panicking handler should yield a panic error, got %v", err)
	}
}
