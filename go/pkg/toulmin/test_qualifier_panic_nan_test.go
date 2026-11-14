//ff:func feature=engine type=engine control=sequence
//ff:what TestQualifierPanicNaN — tests that Qualifier NaN panics
package toulmin

import (
	"math"
	"strings"
	"testing"
)

func TestQualifierPanicNaN(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for NaN qualifier")
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
	g.Rule(WarrantA).Qualifier(math.NaN())
}
