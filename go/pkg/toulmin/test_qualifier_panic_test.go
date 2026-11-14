//ff:func feature=engine type=engine control=sequence
//ff:what TestQualifierPanicAbove — tests that Qualifier above 1.0 panics
package toulmin

import (
	"strings"
	"testing"
)

func TestQualifierPanicAbove(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for qualifier > 1.0")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T", r)
		}
		if !strings.Contains(msg, "qualifier must be between 0.0 and 1.0") {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()
	g := NewGraph("test")
	g.Rule(WarrantA).Qualifier(1.5)
}
