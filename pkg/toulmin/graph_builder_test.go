package toulmin

import (
	"math"
	"testing"
)

func WarrantA(claim any, ground any) bool { return true }
func RebuttalB(claim any, ground any) bool { return true }
func DefeaterC(claim any, ground any) bool { return true }

func TestGraphBuilderWarrantOnly(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, 1.0)
	results := g.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}

func TestGraphBuilderWithDefeat(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, 1.0).
		Rebuttal(RebuttalB, 1.0).
		Defeat(RebuttalB, WarrantA)
	results := g.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}

func TestGraphBuilderCompensation(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, 1.0).
		Rebuttal(RebuttalB, 1.0).
		Defeater(DefeaterC, 1.0).
		Defeat(RebuttalB, WarrantA).
		Defeat(DefeaterC, RebuttalB)
	results := g.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}

func TestGraphBuilderQualifierDefault(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA)
	results := g.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (default qualifier), got %f", results[0].Verdict)
	}
}

func TestGraphBuilderFuncReuse(t *testing.T) {
	g1 := NewGraph("graph1").
		Warrant(WarrantA, 1.0).
		Rebuttal(RebuttalB, 1.0).
		Defeat(RebuttalB, WarrantA)
	g2 := NewGraph("graph2").
		Warrant(WarrantA, 1.0).
		Defeater(DefeaterC, 1.0).
		Defeat(DefeaterC, WarrantA)
	r1 := g1.Evaluate(nil, nil)
	r2 := g2.Evaluate(nil, nil)
	if len(r1) != 1 || len(r2) != 1 {
		t.Fatalf("expected 1 result each, got %d and %d", len(r1), len(r2))
	}
	if r1[0].Verdict != 0.0 {
		t.Errorf("g1: expected 0.0, got %f", r1[0].Verdict)
	}
	if r2[0].Verdict != 0.0 {
		t.Errorf("g2: expected 0.0, got %f", r2[0].Verdict)
	}
}

func TestFuncName(t *testing.T) {
	name := FuncName(WarrantA)
	if name != "WarrantA" {
		t.Errorf("expected 'WarrantA', got '%s'", name)
	}
}
