package toulmin

import (
	"math"
	"testing"
)

func TestWarrantOnly(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}

func TestWarrantWithDefeater(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}

func TestCompensation(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D1", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D2", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"D1"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// D2 attacks D1 → D1 weakened → W partially restored
	// raw(D2)=1.0, raw(D1)=0.5, raw(W)=0.667 → verdict(W)≈+0.333
	expected := 1.0/3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}

func TestStrictWarrant(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Strict,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (strict rejects attack), got %f", results[0].Verdict)
	}
}

func TestCircularAttack(t *testing.T) {
	graph := &RuleGraph{
		Nodes: map[string]*Node{
			"A": {Name: "A", Qualifier: 1.0, Strength: Defeasible},
			"B": {Name: "B", Qualifier: 1.0, Strength: Defeasible},
		},
		Edges: map[string][]string{
			"A": {"B"},
			"B": {"A"},
		},
	}
	verdict := CalcAcceptability("A", graph, 0)
	if verdict < -1.0 || verdict > 1.0 {
		t.Errorf("verdict out of range [-1,1]: %f", verdict)
	}
}

func TestParseAnnotation(t *testing.T) {
	lines := []string{
		`//tm:backing "Böhm-Jacopini theorem"`,
	}
	meta := ParseAnnotation(lines)
	if meta.Backing != "Böhm-Jacopini theorem" {
		t.Errorf("backing: expected 'Böhm-Jacopini theorem', got '%s'", meta.Backing)
	}
	if meta.Qualifier != 1.0 {
		t.Errorf("qualifier: expected default 1.0, got %f", meta.Qualifier)
	}
}
