package toulmin

import (
	"math"
	"testing"
)

func WarrantA(claim any, ground any, backing any) (bool, any)  { return true, nil }
func RebuttalB(claim any, ground any, backing any) (bool, any) { return true, nil }
func DefeaterC(claim any, ground any, backing any) (bool, any) { return true, nil }
func InactiveR(claim any, ground any, backing any) (bool, any) { return false, nil }

func TestGraphWarrantOnly(t *testing.T) {
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 1.0)
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

func TestGraphWithDefeat(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	g.Defeat(r, w)
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

func TestGraphCompensation(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	d := g.Defeater(DefeaterC, nil, 1.0)
	g.Defeat(r, w)
	g.Defeat(d, r)
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

func TestGraphFuncReuse(t *testing.T) {
	g1 := NewGraph("graph1")
	w1 := g1.Warrant(WarrantA, nil, 1.0)
	r1 := g1.Rebuttal(RebuttalB, nil, 1.0)
	g1.Defeat(r1, w1)

	g2 := NewGraph("graph2")
	w2 := g2.Warrant(WarrantA, nil, 1.0)
	d2 := g2.Defeater(DefeaterC, nil, 1.0)
	g2.Defeat(d2, w2)

	res1, err := g1.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("g1 error: %v", err)
	}
	res2, err := g2.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("g2 error: %v", err)
	}
	if len(res1) != 1 || len(res2) != 1 {
		t.Fatalf("expected 1 result each, got %d and %d", len(res1), len(res2))
	}
	if res1[0].Verdict != 0.0 {
		t.Errorf("g1: expected 0.0, got %f", res1[0].Verdict)
	}
	if res2[0].Verdict != 0.0 {
		t.Errorf("g2: expected 0.0, got %f", res2[0].Verdict)
	}
}

func TestGraphTraceAllRules(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 0.8)
	g.Defeat(r, w)
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

func TestGraphTraceIncludesInactive(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(InactiveR, nil, 1.0)
	g.Defeat(r, w)
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
	g := NewGraph("test")
	w := g.Warrant(falseWarrant, nil, 1.0)
	r := g.Rebuttal(trackedRebuttal, nil, 1.0)
	g.Defeat(r, w)
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
	g := NewGraph("test")
	wA := g.Warrant(WarrantA, nil, 1.0)
	wX := g.Warrant(warrantX, nil, 1.0)
	ud := g.Defeater(unrelatedDefeater, nil, 1.0)
	rB := g.Rebuttal(RebuttalB, nil, 1.0)
	g.Defeat(rB, wA)
	g.Defeat(ud, wX)
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
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 0.0)
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

func TestGraphCycleError(t *testing.T) {
	cycleA := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	cycleB := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	a := g.Warrant(cycleA, nil, 1.0)
	b := g.Rebuttal(cycleB, nil, 1.0)
	g.Defeat(b, a)
	g.Defeat(a, b)
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
	isInRole := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Warrant(isInRole, "admin", 1.0)
	g.Warrant(isInRole, "editor", 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (same func, different backing), got %d", len(results))
	}
}

func TestBackingInTrace(t *testing.T) {
	isInRole := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Warrant(isInRole, "admin", 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Trace[0].Backing != "admin" {
		t.Errorf("expected backing 'admin', got %v", results[0].Trace[0].Backing)
	}
}

func TestBackingNil(t *testing.T) {
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Trace[0].Backing != nil {
		t.Errorf("expected backing nil, got %v", results[0].Trace[0].Backing)
	}
}

func TestDefeatWithBacking(t *testing.T) {
	isIPInList := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	isAuth := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	auth := g.Warrant(isAuth, nil, 1.0)
	blocked := g.Rebuttal(isIPInList, "blocklist", 1.0)
	allowed := g.Defeater(isIPInList, "whitelist", 1.0)
	g.Defeat(blocked, auth)
	g.Defeat(allowed, blocked)
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
	g := NewGraph("test")
	g.Warrant(legacyFn, nil, 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
