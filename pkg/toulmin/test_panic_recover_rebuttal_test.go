//ff:func feature=engine type=engine control=sequence
//ff:what TestPanicRecoverRebuttal — tests that panicking rebuttal returns error
package toulmin

import (
	"strings"
	"testing"
)

func TestPanicRecoverRebuttal(t *testing.T) {
	panicRebuttal := func(ctx Context, specs Specs) (bool, any) {
		panic("unexpected")
	}
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(panicRebuttal)
	r.Attacks(w)
	results, err := g.Evaluate(nil)
	if err == nil {
		t.Fatalf("expected error from panicking rebuttal, got results: %+v", results)
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Errorf("expected panic error message, got: %v", err)
	}
}
