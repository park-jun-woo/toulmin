//ff:func feature=engine type=engine control=sequence
//ff:what TestRunTraceGradient — a handler branches on the continuous verdict read from trace
package toulmin

import "testing"

func TestRunTraceGradient(t *testing.T) {
	var branch string
	decide := func(self TraceEntry, t Trace) error {
		if self.Verdict >= 0.5 {
			branch = "strong"
		} else {
			branch = "weak"
		}
		return nil
	}
	g := NewGraph("gradient")
	g.Rule(WarrantA).Qualifier(0.75).RunOn(decide)

	var captured float64
	read := func(self TraceEntry, t Trace) error {
		captured = self.Verdict
		return nil
	}
	g2 := NewGraph("gradient-weak")
	g2.Rule(WarrantA).Qualifier(0.6).RunOn(read)

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
