//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextRecordTrace — tests evalContext.recordTrace for duration on/off, success/panic, and explicit/inferred role branches
package toulmin

import "testing"

func TestEvalContextRecordTrace(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"DurationSuccessExplicitRole", func(t *testing.T) {
			activeFn := func(ctx Context, specs Specs) (bool, any) { return true, "ev" }
			ec := &evalContext{
				fnMap:       map[string]func(Context, Specs) (bool, any){"a": activeFn},
				ran:         map[string]bool{},
				active:      map[string]bool{},
				evidence:    map[string]any{},
				specsMap:    map[string]Specs{},
				strMap:      map[string]Strength{"a": Strict},
				attackerSet: map[string]bool{},
				roleMap:     map[string]string{"a": "rule"},
				qualMap:     map[string]float64{"a": 1.0},
			}
			ec.recordTrace("a", NewContext(), true)
			if ec.err != nil {
				t.Fatalf("unexpected error: %v", ec.err)
			}
			if !ec.ran["a"] {
				t.Fatal("expected ran[a] = true")
			}
			if !ec.active["a"] {
				t.Fatal("expected active[a] = true")
			}
			if len(ec.trace) != 1 {
				t.Fatalf("expected 1 trace entry, got %d", len(ec.trace))
			}
			if ec.trace[0].Role != "rule" {
				t.Errorf("expected explicit role %q, got %q", "rule", ec.trace[0].Role)
			}
			if ec.trace[0].Evidence != "ev" {
				t.Errorf("expected evidence %q, got %v", "ev", ec.trace[0].Evidence)
			}
		}},
		{"DurationPanicSetsErr", func(t *testing.T) {
			panicFn := func(ctx Context, specs Specs) (bool, any) { panic("boom") }
			ec := &evalContext{
				fnMap:       map[string]func(Context, Specs) (bool, any){"a": panicFn},
				ran:         map[string]bool{},
				active:      map[string]bool{},
				evidence:    map[string]any{},
				specsMap:    map[string]Specs{},
				strMap:      map[string]Strength{},
				attackerSet: map[string]bool{},
				roleMap:     map[string]string{},
			}
			ec.recordTrace("a", NewContext(), true)
			if ec.err == nil {
				t.Fatal("expected ec.err to be set after panic")
			}
			if len(ec.trace) != 0 {
				t.Fatalf("expected no trace entry recorded on error, got %d", len(ec.trace))
			}
			if !ec.ran["a"] {
				t.Fatal("expected ran[a] = true even on panic")
			}
		}},
		{"NoDurationSuccessInferredRole", func(t *testing.T) {
			activeFn := func(ctx Context, specs Specs) (bool, any) { return false, nil }
			ec := &evalContext{
				fnMap:       map[string]func(Context, Specs) (bool, any){"a": activeFn},
				ran:         map[string]bool{},
				active:      map[string]bool{},
				evidence:    map[string]any{},
				specsMap:    map[string]Specs{},
				strMap:      map[string]Strength{"a": Defeasible},
				attackerSet: map[string]bool{"a": true},
				roleMap:     map[string]string{},
				qualMap:     map[string]float64{"a": 0.5},
			}
			ec.recordTrace("a", NewContext(), false)
			if ec.err != nil {
				t.Fatalf("unexpected error: %v", ec.err)
			}
			if ec.active["a"] {
				t.Fatal("expected active[a] = false")
			}
			if len(ec.trace) != 1 {
				t.Fatalf("expected 1 trace entry, got %d", len(ec.trace))
			}
			if ec.trace[0].Role == "" {
				t.Error("expected inferred non-empty role")
			}
			if ec.trace[0].Duration != 0 {
				t.Errorf("expected zero duration when duration=false, got %v", ec.trace[0].Duration)
			}
		}},
		{"NoDurationPanicSetsErr", func(t *testing.T) {
			panicFn := func(ctx Context, specs Specs) (bool, any) { panic("boom") }
			ec := &evalContext{
				fnMap:       map[string]func(Context, Specs) (bool, any){"a": panicFn},
				ran:         map[string]bool{},
				active:      map[string]bool{},
				evidence:    map[string]any{},
				specsMap:    map[string]Specs{},
				strMap:      map[string]Strength{},
				attackerSet: map[string]bool{},
				roleMap:     map[string]string{},
			}
			ec.recordTrace("a", NewContext(), false)
			if ec.err == nil {
				t.Fatal("expected ec.err to be set after panic")
			}
			if len(ec.trace) != 0 {
				t.Fatalf("expected no trace entry recorded on error, got %d", len(ec.trace))
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
