//ff:func feature=engine type=engine control=sequence
//ff:what TestEvalContextEdgeForce — edgeForce returns the stored force, or 1 when none is stored
package toulmin

import "testing"

func TestEvalContextEdgeForce(t *testing.T) {
	ec := &evalContext{forces: map[string][]float64{"a": {0.25, 0.0}}}
	if f := ec.edgeForce("a", 0); f != 0.25 {
		t.Errorf("expected 0.25, got %v", f)
	}
	if f := ec.edgeForce("a", 1); f != 0.0 {
		t.Errorf("expected stored 0.0 (not the default), got %v", f)
	}
	if f := ec.edgeForce("a", 2); f != 1.0 {
		t.Errorf("expected default 1.0 past stored forces, got %v", f)
	}
	if f := ec.edgeForce("b", 0); f != 1.0 {
		t.Errorf("expected default 1.0 for edges without stored forces, got %v", f)
	}
}
