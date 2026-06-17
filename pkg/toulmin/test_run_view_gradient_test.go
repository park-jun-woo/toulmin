//ff:func feature=engine type=engine control=sequence
//ff:what TestRunViewGradient — a handler branches on the continuous verdict read via view.Get
package toulmin

import "testing"

func TestRunViewGradient(t *testing.T) {
	var branch string
	decide := func(ctx Context, ev NodeEvent, view RunView) error {
		v, _ := view.Get("WarrantA")
		if v.Verdict >= 0.5 {
			branch = "strong"
		} else {
			branch = "weak"
		}
		return nil
	}
	g := NewGraph("gradient")
	g.Rule(WarrantA).Qualifier(0.75).OnActive(decide).OnDefeated(decide).OnInactive(decide)

	var captured float64
	read := func(ctx Context, ev NodeEvent, view RunView) error {
		v, _ := view.Get("WarrantA")
		captured = v.Verdict
		return nil
	}
	g2 := NewGraph("gradient-weak")
	g2.Rule(WarrantA).Qualifier(0.6).OnActive(read).OnDefeated(read).OnInactive(read)

	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	// Qualifier 0.75, unattacked → verdict 2*0.75-1 = 0.5 → strong branch.
	if branch != "strong" {
		t.Errorf("verdict 0.5 should take the strong branch, got %q", branch)
	}

	if _, _, err := g2.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	// Qualifier 0.6, unattacked → verdict 2*0.6-1 = 0.2 (< 0.5).
	if captured >= 0.5 {
		t.Errorf("verdict should be below the 0.5 threshold, got %f", captured)
	}
}
