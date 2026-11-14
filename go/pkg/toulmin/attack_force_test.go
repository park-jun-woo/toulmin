//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestAttackForce — tests attackForce default, valid values, and panics on extra values, NaN, and out-of-range
package toulmin

import (
	"math"
	"strings"
	"testing"
)

func TestAttackForce(t *testing.T) {
	expectPanic := func(t *testing.T, force []float64, want string) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("force %v: expected panic", force)
			}
			if msg, _ := r.(string); !strings.Contains(msg, want) {
				t.Fatalf("force %v: unexpected panic %v", force, r)
			}
		}()
		attackForce(force)
	}
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"DefaultIsOne", func(t *testing.T) {
			if f := attackForce(nil); f != 1.0 {
				t.Fatalf("expected default 1.0, got %v", f)
			}
		}},
		{"ValidValues", func(t *testing.T) {
			for _, want := range []float64{0.0, 0.3, 1.0} {
				if f := attackForce([]float64{want}); f != want {
					t.Errorf("expected %v, got %v", want, f)
				}
			}
		}},
		{"TooMany", func(t *testing.T) { expectPanic(t, []float64{0.5, 0.5}, "at most one force") }},
		{"NaN", func(t *testing.T) { expectPanic(t, []float64{math.NaN()}, "between 0.0 and 1.0") }},
		{"Below", func(t *testing.T) { expectPanic(t, []float64{-0.1}, "between 0.0 and 1.0") }},
		{"Above", func(t *testing.T) { expectPanic(t, []float64{1.1}, "between 0.0 and 1.0") }},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
