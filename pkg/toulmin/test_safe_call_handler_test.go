//ff:func feature=engine type=engine control=sequence
//ff:what TestSafeCallHandler — safeCallHandler returns nil, the handler error, or a recovered panic
package toulmin

import (
	"fmt"
	"strings"
	"testing"
)

func TestSafeCallHandler(t *testing.T) {
	ev := NodeEvent{Name: "n"}

	if err := safeCallHandler(func(ctx Context, ev NodeEvent, view RunView) error { return nil }, NewContext(), ev, nil); err != nil {
		t.Errorf("ok handler should return nil, got %v", err)
	}

	want := fmt.Errorf("bad")
	if err := safeCallHandler(func(ctx Context, ev NodeEvent, view RunView) error { return want }, NewContext(), ev, nil); err != want {
		t.Errorf("error handler should propagate error, got %v", err)
	}

	err := safeCallHandler(func(ctx Context, ev NodeEvent, view RunView) error { panic("boom") }, NewContext(), ev, nil)
	if err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Errorf("panicking handler should yield a panic error, got %v", err)
	}
}
