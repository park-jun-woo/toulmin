//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextEvalRule — tests evalContext.evalRule for Trace-enabled and default (calc) branches
package toulmin

import "testing"

func TestEvalContextEvalRule(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Trace", func(t *testing.T) {
			activeFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){"a": activeFn},
				verdictCache: map[string]float64{},
				ran:          map[string]bool{},
				active:       map[string]bool{},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
				strMap:       map[string]Strength{"a": Strict},
				attackerSet:  map[string]bool{},
				roleMap:      map[string]string{},
				qualMap:      map[string]float64{"a": 1.0},
			}
			v := ec.evalRule("a", NewContext(), EvalOption{Trace: true})
			if v != 1.0 {
				t.Fatalf("expected 1.0, got %v", v)
			}
			if len(ec.trace) != 1 {
				t.Fatalf("expected trace to be recorded, got %d entries", len(ec.trace))
			}
		}},
		{"Default", func(t *testing.T) {
			activeFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){"a": activeFn},
				verdictCache: map[string]float64{},
				ran:          map[string]bool{},
				active:       map[string]bool{},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
				strMap:       map[string]Strength{"a": Strict},
				qualMap:      map[string]float64{"a": 1.0},
			}
			v := ec.evalRule("a", NewContext(), EvalOption{Trace: false})
			if v != 1.0 {
				t.Fatalf("expected 1.0, got %v", v)
			}
			if len(ec.trace) != 0 {
				t.Fatalf("expected no trace recorded for non-trace path, got %d entries", len(ec.trace))
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
