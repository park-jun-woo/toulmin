//ff:func feature=engine type=engine control=sequence
//ff:what TestRecursiveRejected — tests that Recursive eval method returns not-yet-implemented error
package toulmin

import (
	"strings"
	"testing"
)

func TestRecursiveRejected(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA)
	_, err := g.Evaluate(NewContext(), EvalOption{Method: Recursive})
	if err == nil {
		t.Fatal("expected error for Recursive method")
	}
	if !strings.Contains(err.Error(), "not yet implemented") {
		t.Errorf("expected 'not yet implemented' in error, got: %v", err)
	}
}
