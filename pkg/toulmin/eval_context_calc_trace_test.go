//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextCalcTrace — tests evalContext.calcTrace's branches: prior error, cache hit, missing/nil fn, already-ran skip, record-trace-panic error, inactive, Strict vs Defeasible attacker loop
package toulmin

import (
	"errors"
	"testing"
)

func TestEvalContextCalcTrace(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"PriorError", func(t *testing.T) {
			ec := &evalContext{err: errors.New("boom")}
			if v := ec.calcTrace("a", NewContext(), false); v != -1.0 {
				t.Fatalf("expected -1.0, got %v", v)
			}
		}},
		{"CacheHit", func(t *testing.T) {
			ec := &evalContext{verdictCache: map[string]float64{"a": 0.42}}
			if v := ec.calcTrace("a", NewContext(), false); v != 0.42 {
				t.Fatalf("expected 0.42, got %v", v)
			}
		}},
		{"MissingFn", func(t *testing.T) {
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){},
				verdictCache: map[string]float64{},
			}
			if v := ec.calcTrace("missing", NewContext(), false); v != -1.0 {
				t.Fatalf("expected -1.0, got %v", v)
			}
		}},
		{"NilFn", func(t *testing.T) {
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){"a": nil},
				verdictCache: map[string]float64{},
			}
			if v := ec.calcTrace("a", NewContext(), false); v != -1.0 {
				t.Fatalf("expected -1.0, got %v", v)
			}
		}},
		{"RecordTracePanicSetsErr", func(t *testing.T) {
			panicFn := func(ctx Context, specs Specs) (bool, any) { panic("kaboom") }
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){"a": panicFn},
				verdictCache: map[string]float64{},
				ran:          map[string]bool{},
				active:       map[string]bool{},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
				strMap:       map[string]Strength{},
				attackerSet:  map[string]bool{},
				roleMap:      map[string]string{},
			}
			if v := ec.calcTrace("a", NewContext(), false); v != -1.0 {
				t.Fatalf("expected -1.0, got %v", v)
			}
			if ec.err == nil {
				t.Fatal("expected ec.err to be set after panic")
			}
		}},
		{"AlreadyRanInactive", func(t *testing.T) {
			calledAgain := false
			fn := func(ctx Context, specs Specs) (bool, any) { calledAgain = true; return true, nil }
			ec := &evalContext{
				fnMap:        map[string]func(Context, Specs) (bool, any){"a": fn},
				verdictCache: map[string]float64{},
				ran:          map[string]bool{"a": true},
				active:       map[string]bool{"a": false},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
				strMap:       map[string]Strength{},
				attackerSet:  map[string]bool{},
				roleMap:      map[string]string{},
			}
			if v := ec.calcTrace("a", NewContext(), false); v != -1.0 {
				t.Fatalf("expected -1.0 for inactive node, got %v", v)
			}
			if calledAgain {
				t.Fatal("expected fn not to be re-invoked for already-ran node")
			}
		}},
		{"StrictSkipsAttackerLoop", func(t *testing.T) {
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
				// Present edges must be ignored because strength is Strict.
				edges:   map[string][]string{"a": {"unreachable"}},
				qualMap: map[string]float64{"a": 1.0},
			}
			v := ec.calcTrace("a", NewContext(), false)
			if v != 1.0 {
				t.Fatalf("expected 1.0, got %v", v)
			}
		}},
		{"DefeasibleWithAttackerLoop", func(t *testing.T) {
			activeFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			ec := &evalContext{
				fnMap: map[string]func(Context, Specs) (bool, any){
					"b":  activeFn,
					"a1": activeFn,
					"a2": activeFn,
				},
				verdictCache: map[string]float64{},
				ran:          map[string]bool{},
				active:       map[string]bool{},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
				attackerSet:  map[string]bool{"a1": true, "a2": true},
				roleMap:      map[string]string{},
				strMap: map[string]Strength{
					"b":  Defeasible,
					"a1": Strict,
					"a2": Strict,
				},
				edges: map[string][]string{
					"b": {"a1", "a2"},
				},
				qualMap: map[string]float64{
					"b":  1.0,
					"a1": 1.0,
					"a2": 0.5,
				},
			}
			// a1 -> v=1.0 (raw=1.0), a2 -> v=0.0 (raw=0.5).
			// sum = (1.0+1.0)/2 + (0.0+1.0)/2 = 1.0 + 0.5 = 1.5
			// raw = qual[b]/(1+sum) = 1/2.5 = 0.4, v = 2*0.4-1 = -0.2
			v := ec.calcTrace("b", NewContext(), true)
			if v < -0.2000001 || v > -0.1999999 {
				t.Fatalf("expected -0.2, got %v", v)
			}
			if ec.verdictCache["b"] != v {
				t.Fatalf("expected verdictCache to be populated with %v, got %v", v, ec.verdictCache["b"])
			}
			if len(ec.trace) != 3 {
				t.Fatalf("expected 3 trace entries recorded, got %d", len(ec.trace))
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
