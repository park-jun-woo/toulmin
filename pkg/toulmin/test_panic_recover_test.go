//ff:func feature=engine type=engine control=sequence
//ff:what TestPanicRecover — tests that panicking warrant returns error instead of crashing
package toulmin

import (
	"strings"
	"testing"
)

func TestPanicRecover(t *testing.T) {
	panicRule := func(claim, ground, backing any) (bool, any) {
		_ = ground.(string)
		return true, nil
	}
	g := NewGraph("test")
	g.Warrant(panicRule, nil, 1.0)
	results, err := g.Evaluate(nil, nil)
	if err == nil {
		t.Fatalf("expected error from panicking rule, got results: %+v", results)
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Errorf("expected panic error message, got: %v", err)
	}
}
