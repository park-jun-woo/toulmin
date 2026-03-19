package toulmin

import (
	"math"
	"testing"
)

func WarrantA(claim any, ground any, backing any) (bool, any)  { return true, nil }
func RebuttalB(claim any, ground any, backing any) (bool, any) { return true, nil }
func DefeaterC(claim any, ground any, backing any) (bool, any) { return true, nil }
func InactiveR(claim any, ground any, backing any) (bool, any) { return false, nil }

func TestGraphBuilderWarrantOnly(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}

func TestGraphBuilderWithDefeat(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0).
		Rebuttal(RebuttalB, nil, 1.0).
		Defeat(RebuttalB, WarrantA)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}

func TestGraphBuilderCompensation(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0).
		Rebuttal(RebuttalB, nil, 1.0).
		Defeater(DefeaterC, nil, 1.0).
		Defeat(RebuttalB, WarrantA).
		Defeat(DefeaterC, RebuttalB)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}

func TestGraphBuilderFuncReuse(t *testing.T) {
	g1 := NewGraph("graph1").
		Warrant(WarrantA, nil, 1.0).
		Rebuttal(RebuttalB, nil, 1.0).
		Defeat(RebuttalB, WarrantA)
	g2 := NewGraph("graph2").
		Warrant(WarrantA, nil, 1.0).
		Defeater(DefeaterC, nil, 1.0).
		Defeat(DefeaterC, WarrantA)
	r1, err := g1.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("g1 error: %v", err)
	}
	r2, err := g2.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("g2 error: %v", err)
	}
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

func TestGraphBuilderTraceAllRules(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0).
		Rebuttal(RebuttalB, nil, 0.8).
		Defeat(RebuttalB, WarrantA)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	trace := results[0].Trace
	if len(trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(trace))
	}
	if trace[0].Name != "WarrantA" || trace[0].Role != "warrant" || !trace[0].Activated || trace[0].Qualifier != 1.0 {
		t.Errorf("trace[0] unexpected: %+v", trace[0])
	}
	if trace[1].Name != "RebuttalB" || trace[1].Role != "rebuttal" || !trace[1].Activated || trace[1].Qualifier != 0.8 {
		t.Errorf("trace[1] unexpected: %+v", trace[1])
	}
}

func TestGraphBuilderTraceIncludesInactive(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0).
		Rebuttal(InactiveR, nil, 1.0).
		Defeat(InactiveR, WarrantA)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	trace := results[0].Trace
	if len(trace) != 2 {
		t.Fatalf("expected 2 trace entries, got %d", len(trace))
	}
	if trace[1].Activated {
		t.Errorf("expected InactiveR activated=false, got true")
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (rebuttal inactive), got %f", results[0].Verdict)
	}
}

func TestLazySkipsRebuttalWhenWarrantFalse(t *testing.T) {
	rebuttalCalled := false
	falseWarrant := func(claim any, ground any, backing any) (bool, any) { return false, nil }
	trackedRebuttal := func(claim any, ground any, backing any) (bool, any) {
		rebuttalCalled = true
		return true, nil
	}
	g := NewGraph("test").
		Warrant(falseWarrant, nil, 1.0).
		Rebuttal(trackedRebuttal, nil, 1.0).
		Defeat(trackedRebuttal, falseWarrant)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results (warrant false), got %d", len(results))
	}
	if rebuttalCalled {
		t.Error("rebuttal func should not be called when warrant is false")
	}
}

func TestTraceOnlyRelevantRules(t *testing.T) {
	warrantX := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	unrelatedDefeater := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0).
		Warrant(warrantX, nil, 1.0).
		Defeater(unrelatedDefeater, nil, 1.0).
		Rebuttal(RebuttalB, nil, 1.0).
		Defeat(RebuttalB, WarrantA).
		Defeat(unrelatedDefeater, warrantX)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	traceA := results[0].Trace
	for _, te := range traceA {
		if te.Name != "WarrantA" && te.Name != "RebuttalB" {
			t.Errorf("WarrantA trace contains unrelated rule: %s", te.Name)
		}
	}
}

func TestFuncName(t *testing.T) {
	name := FuncName(WarrantA)
	if name != "WarrantA" {
		t.Errorf("expected 'WarrantA', got '%s'", name)
	}
}

func TestFuncIDUniqueness(t *testing.T) {
	fn1 := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	fn2 := func(claim any, ground any, backing any) (bool, any) { return false, nil }
	id1 := funcID(fn1)
	id2 := funcID(fn2)
	if id1 == id2 {
		t.Errorf("expected distinct funcIDs for different closures, both got %s", id1)
	}
}

func TestQualifierZero(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 0.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != -1.0 {
		t.Errorf("expected -1.0 (qualifier=0), got %f", results[0].Verdict)
	}
}

func TestDeepDefeatChain(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }
	eng := NewEngine()
	eng.Register(RuleMeta{Name: "W", Qualifier: 1.0, Strength: Defeasible, Fn: fn})
	prev := "W"
	for i := 1; i <= 11; i++ {
		name := "D" + string(rune('0'+i))
		eng.Register(RuleMeta{Name: name, Qualifier: 1.0, Strength: Defeater, Defeats: []string{prev}, Fn: fn})
		prev = name
	}
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := results[0].Verdict
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("verdict should be finite, got %f", v)
	}
}

func TestGraphBuilderCycleError(t *testing.T) {
	cycleA := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	cycleB := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test").
		Warrant(cycleA, nil, 1.0).
		Rebuttal(cycleB, nil, 1.0).
		Defeat(cycleB, cycleA).
		Defeat(cycleA, cycleB)
	_, err := g.Evaluate(nil, nil)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
}

func TestDeepDefeatChainOver100(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }
	eng := NewEngine()
	eng.Register(RuleMeta{Name: "W", Qualifier: 1.0, Strength: Defeasible, Fn: fn})
	prev := "W"
	for i := 1; i <= 150; i++ {
		name := string(rune('A'+i%26)) + string(rune('0'+i/26))
		eng.Register(RuleMeta{Name: name, Qualifier: 1.0, Strength: Defeater, Defeats: []string{prev}, Fn: fn})
		prev = name
	}
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := results[0].Verdict
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("verdict should be finite for deep non-cyclic chain, got %f", v)
	}
}

func TestBackingSameFunc(t *testing.T) {
	isInRole := func(claim any, ground any, backing any) (bool, any) {
		return true, nil
	}
	g := NewGraph("test").
		Warrant(isInRole, "admin", 1.0).
		Warrant(isInRole, "editor", 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (same func, different backing), got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("admin: expected +1.0, got %f", results[0].Verdict)
	}
	if results[1].Verdict != 1.0 {
		t.Errorf("editor: expected +1.0, got %f", results[1].Verdict)
	}
}

func TestBackingInTrace(t *testing.T) {
	isInRole := func(claim any, ground any, backing any) (bool, any) {
		return true, nil
	}
	g := NewGraph("test").
		Warrant(isInRole, "admin", 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Trace) != 1 {
		t.Fatalf("expected 1 trace entry, got %d", len(results[0].Trace))
	}
	if results[0].Trace[0].Backing != "admin" {
		t.Errorf("expected backing 'admin', got %v", results[0].Trace[0].Backing)
	}
}

func TestBackingNil(t *testing.T) {
	g := NewGraph("test").
		Warrant(WarrantA, nil, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Trace[0].Backing != nil {
		t.Errorf("expected backing nil, got %v", results[0].Trace[0].Backing)
	}
}

func TestDefeatWithBacking(t *testing.T) {
	isIPInList := func(claim any, ground any, backing any) (bool, any) {
		return true, nil
	}
	isAuth := func(claim any, ground any, backing any) (bool, any) {
		return true, nil
	}
	blocklist := "blocklist"
	whitelist := "whitelist"
	g := NewGraph("test").
		Warrant(isAuth, nil, 1.0).
		Rebuttal(isIPInList, blocklist, 1.0).
		Defeater(isIPInList, whitelist, 1.0).
		DefeatWith(isIPInList, blocklist, isAuth, nil).
		DefeatWith(isIPInList, whitelist, isIPInList, blocklist)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f (whitelist defeats blocklist), got %f", expected, results[0].Verdict)
	}
}

func TestLegacySignature(t *testing.T) {
	legacyFn := func(claim any, ground any) (bool, any) { return true, nil }
	g := NewGraph("test").
		Warrant(legacyFn, nil, 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
