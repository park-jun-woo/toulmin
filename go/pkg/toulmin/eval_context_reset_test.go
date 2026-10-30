//ff:func feature=engine type=engine control=sequence
//ff:what TestEvalContextReset — tests evalContext.reset clears ran/active/ruleVerdict/applicability/evidence/trace/err/verdictCache
package toulmin

import (
	"errors"
	"testing"
)

func TestEvalContextReset(t *testing.T) {
	ec := &evalContext{
		ran:           map[string]bool{"a": true},
		active:        map[string]bool{"a": true},
		ruleVerdict:   map[string]float64{"a": 0.5},
		applicability: map[string]float64{"a": 0.75},
		evidence:      map[string]any{"a": "ev"},
		trace:         []TraceEntry{{Name: "a"}},
		err:           errors.New("boom"),
		verdictCache:  map[string]float64{"a": 1.0},
	}
	ec.reset()
	if len(ec.ran) != 0 {
		t.Errorf("expected ran cleared, got %v", ec.ran)
	}
	if len(ec.active) != 0 {
		t.Errorf("expected active cleared, got %v", ec.active)
	}
	if len(ec.ruleVerdict) != 0 || len(ec.applicability) != 0 {
		t.Errorf("expected ruleVerdict/applicability cleared, got %v / %v", ec.ruleVerdict, ec.applicability)
	}
	if len(ec.evidence) != 0 {
		t.Errorf("expected evidence cleared, got %v", ec.evidence)
	}
	if ec.trace != nil {
		t.Errorf("expected trace nil, got %v", ec.trace)
	}
	if ec.err != nil {
		t.Errorf("expected err nil, got %v", ec.err)
	}
	if len(ec.verdictCache) != 0 {
		t.Errorf("expected verdictCache cleared, got %v", ec.verdictCache)
	}
}
