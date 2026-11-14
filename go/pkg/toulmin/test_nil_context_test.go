//ff:func feature=engine type=engine control=sequence
//ff:what TestNilContext — tests that Evaluate with nil context returns error
package toulmin

import (
	"strings"
	"testing"
)

func TestNilContext(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA)
	_, err := g.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for nil context")
	}
	if !strings.Contains(err.Error(), "ctx must not be nil") {
		t.Errorf("expected 'ctx must not be nil' in error, got: %v", err)
	}
}
