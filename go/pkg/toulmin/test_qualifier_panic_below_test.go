//ff:func feature=engine type=engine control=sequence
//ff:what TestQualifierPanicBelow — tests that Qualifier below 0.0 panics
package toulmin

import (
	"strings"
	"testing"
)

func TestQualifierPanicBelow(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for qualifier < 0.0")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T", r)
		}
		if !strings.Contains(msg, "qualifier must be between 0.0 and 1.0") {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()
	fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(fn).Qualifier(-0.1)
}
