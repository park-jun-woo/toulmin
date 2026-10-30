//ff:func feature=engine type=engine control=sequence
//ff:what TestEvalContextSettle — settle computes 2·q·aEff/(1+sum) − 1, records applicability, and caches the verdict
package toulmin

import (
	"math"
	"testing"
)

func TestEvalContextSettle(t *testing.T) {
	ec := &evalContext{
		qualMap:       map[string]float64{"a": 0.5},
		applicability: map[string]float64{},
		verdictCache:  map[string]float64{},
	}
	// raw = 0.5·0.5/(1+1) = 0.125 → v = -0.75
	v := ec.settle("a", 0.5, 1.0)
	if math.Abs(v-(-0.75)) > 1e-12 {
		t.Fatalf("expected -0.75, got %v", v)
	}
	if ec.applicability["a"] != 0.5 {
		t.Errorf("expected applicability 0.5, got %v", ec.applicability["a"])
	}
	if ec.verdictCache["a"] != v {
		t.Errorf("expected cached verdict %v, got %v", v, ec.verdictCache["a"])
	}
}
