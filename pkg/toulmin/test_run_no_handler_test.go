//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunNoHandler — nodes without handlers pass through; Run results match Evaluate
package toulmin

import "testing"

func TestRunNoHandler(t *testing.T) {
	g := NewGraph("nohandler")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.Attacks(w)

	evalResults, err1 := g.Evaluate(NewContext())
	if err1 != nil {
		t.Fatalf("evaluate error: %v", err1)
	}
	runResults, view, err2 := g.Run(NewContext())
	if err2 != nil {
		t.Fatalf("run error: %v", err2)
	}
	if got := len(view.All()); got != 2 {
		t.Errorf("RunView must snapshot every node even with no handlers, want 2, got %d", got)
	}
	if len(runResults) != len(evalResults) {
		t.Fatalf("result count mismatch: run %d vs evaluate %d", len(runResults), len(evalResults))
	}
	for i := range evalResults {
		if runResults[i].Verdict != evalResults[i].Verdict {
			t.Errorf("verdict mismatch at %d: %f vs %f", i, runResults[i].Verdict, evalResults[i].Verdict)
		}
		if runResults[i].Name != evalResults[i].Name {
			t.Errorf("name mismatch at %d: %s vs %s", i, runResults[i].Name, evalResults[i].Name)
		}
	}
}
