//ff:func feature=engine type=engine control=sequence
//ff:what TestWrongSignaturePanic — tests that Rule with wrong func signature panics
package toulmin

import (
	"strings"
	"testing"
)

func TestWrongSignaturePanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for wrong function signature")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T", r)
		}
		if !strings.Contains(msg, "fn must be func(Context, Specs)(bool, any)") {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()
	g := NewGraph("test")
	g.Rule(func() {})
}
