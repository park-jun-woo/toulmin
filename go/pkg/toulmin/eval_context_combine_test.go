//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextCombine — combine: Strict ignores attackers, rebuttals go to the denominator with force, exceptions gate applicability
package toulmin

import (
	"math"
	"testing"
)

func TestEvalContextCombine(t *testing.T) {
	newEC := func(str Strength, edges []string, forces []float64) *evalContext {
		return &evalContext{
			ruleVerdict:   map[string]float64{"w": 0.0}, // a = 0.5
			qualMap:       map[string]float64{"w": 1.0},
			strMap:        map[string]Strength{"w": str, "c": Defeasible, "x": Defeater},
			roleMap:       map[string]string{"c": "counter", "x": "except"},
			edges:         map[string][]string{"w": edges},
			forces:        map[string][]float64{"w": forces},
			applicability: map[string]float64{},
			verdictCache:  map[string]float64{},
		}
	}
	attacker := map[string]float64{"c": 1.0, "x": 0.0} // c fully applies; x half applies
	called := 0
	eval := func(id string) float64 { called++; return attacker[id] }
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-12 }
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"StrictIgnoresAttackers", func(t *testing.T) {
			called = 0
			ec := newEC(Strict, []string{"c", "x"}, nil)
			if v := ec.combine("w", eval); v != 0.0 || called != 0 {
				t.Fatalf("expected 0.0 without evaluating attackers, got %v (called %d)", v, called)
			}
			if ec.applicability["w"] != 0.5 {
				t.Errorf("strict applicability must be ungated a=0.5, got %v", ec.applicability["w"])
			}
		}},
		{"RebuttalWithForce", func(t *testing.T) {
			// S = 0.5·(1+1)/2 = 0.5 → raw = 0.5/1.5 → v = 2/3 − 1 = −1/3
			ec := newEC(Defeasible, []string{"c"}, []float64{0.5})
			if v := ec.combine("w", eval); !near(v, -1.0/3.0) {
				t.Fatalf("expected -1/3, got %v", v)
			}
			if ec.applicability["w"] != 0.5 {
				t.Errorf("rebuttal must not touch applicability, got %v", ec.applicability["w"])
			}
		}},
		{"ExceptGate", func(t *testing.T) {
			// gate = 1 − 1·(0+1)/2 = 0.5 → a_eff = 0.25 → v = −0.5
			ec := newEC(Defeasible, []string{"x"}, nil)
			if v := ec.combine("w", eval); !near(v, -0.5) {
				t.Fatalf("expected -0.5, got %v", v)
			}
			if ec.applicability["w"] != 0.25 {
				t.Errorf("expected applicability 0.25, got %v", ec.applicability["w"])
			}
		}},
		{"Mixed", func(t *testing.T) {
			// a_eff = 0.5·(1 − 0.5·0.5) = 0.375, S = 1 → raw = 0.1875 → v = −0.625
			ec := newEC(Defeasible, []string{"c", "x"}, []float64{1.0, 0.5})
			if v := ec.combine("w", eval); !near(v, -0.625) {
				t.Fatalf("expected -0.625, got %v", v)
			}
			if ec.applicability["w"] != 0.375 {
				t.Errorf("expected applicability 0.375, got %v", ec.applicability["w"])
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
