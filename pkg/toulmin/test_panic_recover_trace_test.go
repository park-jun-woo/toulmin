//ff:func feature=engine type=engine control=sequence
//ff:what TestPanicRecoverTrace — tests that panicking rule returns error in trace mode
package toulmin

import (
	"strings"
	"testing"
)

func TestPanicRecoverTrace(t *testing.T) {
	panicRule := func(claim any, ground any, backing Backing) (bool, any) {
		_ = ground.(string)
		return true, nil
	}
	g := NewGraph("test")
	g.Warrant(panicRule, nil, 1.0)
	results, err := g.Evaluate(nil, nil, EvalOption{Trace: true})
	if err == nil {
		t.Fatalf("expected error from panicking rule, got results: %+v", results)
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Errorf("expected panic error message, got: %v", err)
	}
}
