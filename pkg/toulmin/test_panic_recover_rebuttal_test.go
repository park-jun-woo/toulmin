//ff:func feature=engine type=engine control=sequence
//ff:what TestPanicRecoverRebuttal — tests that panicking rebuttal returns error
package toulmin

import (
	"strings"
	"testing"
)

func TestPanicRecoverRebuttal(t *testing.T) {
	panicRebuttal := func(claim, ground, backing any) (bool, any) {
		panic("unexpected")
	}
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(panicRebuttal, nil, 1.0)
	g.Defeat(r, w)
	results, err := g.Evaluate(nil, nil)
	if err == nil {
		t.Fatalf("expected error from panicking rebuttal, got results: %+v", results)
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Errorf("expected panic error message, got: %v", err)
	}
}
