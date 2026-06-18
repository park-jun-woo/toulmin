//ff:func feature=engine type=engine control=sequence
//ff:what TestSafeCallHandler — safeCallHandler returns nil, the handler error, or a recovered panic
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

func TestSafeCallHandler(t *testing.T) {
	tr := Trace{nodes: []TraceEntry{{Name: "n"}}, ctx: NewContext()}

	if err := safeCallHandler(func(t Trace) error { return nil }, tr); err != nil {
		t.Errorf("ok handler should return nil, got %v", err)
	}

	want := fmt.Errorf("bad")
	if err := safeCallHandler(func(t Trace) error { return want }, tr); err != want {
		t.Errorf("error handler should propagate error, got %v", err)
	}

	err := safeCallHandler(func(t Trace) error { panic("boom") }, tr)
	if err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Errorf("panicking handler should yield a panic error, got %v", err)
	}
}
