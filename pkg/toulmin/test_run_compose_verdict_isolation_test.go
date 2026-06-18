//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeVerdictIsolation — a sub-graph's verdict does not leak into the parent's results
package toulmin

import "testing"

func TestRunComposeVerdictIsolation(t *testing.T) {
	// Active sub-warrant with a sub-unity qualifier: its verdict (0.5 here) must not
	// leak into the parent's result. The RunOn handler captures the sub node's own
	// verdict to prove it is computed independently of the parent.
	subRan := false
	var subVerdict float64
	subWarrant := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	sub := NewGraph("sub")
	sub.Rule(subWarrant).Qualifier(0.75).RunOn(func(ctx Context, self TraceEntry, trace []TraceEntry) error {
		subRan = true
		subVerdict = self.Verdict
		return nil
	})

	parent := NewGraph("parent")
	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent.Rule(active).Run(sub)

	results, _, err := parent.Run(NewContext())
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !subRan {
		t.Fatal("sub-graph did not Run")
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 parent result, got %d", len(results))
	}
	// Parent warrant is unattacked with default qualifier 1.0 → verdict 1.0,
	// regardless of the sub-graph's own verdict.
	if results[0].Verdict != 1.0 {
		t.Errorf("parent verdict must stay 1.0 regardless of sub-graph verdict, got %f", results[0].Verdict)
	}
	// Sub warrant qualifier 0.75, unattacked → verdict 2*0.75-1 = 0.5, independent of parent.
	if subVerdict != 0.5 {
		t.Errorf("sub-graph verdict must be its own (0.5), got %f", subVerdict)
	}
}
