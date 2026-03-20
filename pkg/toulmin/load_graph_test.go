package toulmin

import (
	"testing"
)

func TestLoadGraph_Basic(t *testing.T) {
	wFn := func(c any, g any, b any) (bool, any) { return true, nil }
	rFn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
			{Name: "R", Role: "rebuttal"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "W"},
		},
	}

	funcs := map[string]any{"W": wFn, "R": rFn}

	g, err := LoadGraph(def, funcs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("evaluate error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
}

func TestLoadGraph_WithBacking(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) {
		return b.(string) == "admin", b
	}

	def := GraphDef{
		Graph: "backing-test",
		Rules: []GraphRuleDef{
			{Name: "checkRole", Role: "warrant", Qualifier: 1.0},
		},
	}

	funcs := map[string]any{"checkRole": fn}
	backings := map[string]any{"checkRole": "admin"}

	g, err := LoadGraph(def, funcs, backings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil, nil)
	if len(results) == 0 || results[0].Verdict <= 0 {
		t.Error("expected positive verdict with admin backing")
	}
}

func TestLoadGraph_WithDefeater(t *testing.T) {
	wFn := func(c any, g any, b any) (bool, any) { return true, nil }
	rFn := func(c any, g any, b any) (bool, any) { return true, nil }
	dFn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "defeater-test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
			{Name: "R", Role: "rebuttal"},
			{Name: "D", Role: "defeater"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "W"},
			{From: "D", To: "R"},
		},
	}

	funcs := map[string]any{"W": wFn, "R": rFn, "D": dFn}

	g, err := LoadGraph(def, funcs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil, nil)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	// W is attacked by R, but R is defeated by D → W should prevail
	if results[0].Verdict <= 0 {
		t.Errorf("expected positive verdict (defeater neutralizes rebuttal), got %f", results[0].Verdict)
	}
}

func TestLoadGraph_MissingFunc(t *testing.T) {
	def := GraphDef{
		Graph: "missing",
		Rules: []GraphRuleDef{
			{Name: "unknown", Role: "warrant"},
		},
	}

	_, err := LoadGraph(def, map[string]any{}, nil)
	if err == nil {
		t.Error("expected error for missing function")
	}
}

func TestLoadGraph_InvalidRole(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "bad-role",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "invalid"},
		},
	}

	_, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err == nil {
		t.Error("expected error for invalid role")
	}
}

func TestLoadGraph_MissingDefeatRef(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "bad-edge",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
		},
		Defeats: []GraphEdgeDef{
			{From: "ghost", To: "W"},
		},
	}

	_, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err == nil {
		t.Error("expected error for missing defeat reference")
	}
}

func TestLoadGraph_DefaultQualifier(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "default-q",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"}, // Qualifier=0 → default 1.0
		},
	}

	g, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil, nil)
	if len(results) == 0 || results[0].Verdict != 1.0 {
		t.Errorf("expected verdict 1.0 with default qualifier, got %v", results)
	}
}
