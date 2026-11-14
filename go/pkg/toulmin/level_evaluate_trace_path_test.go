//ff:func feature=engine type=engine control=iteration dimension=2
//ff:what TestLevelEvaluateTracePath — Evaluate+Trace path (recordTrace) fills Level identically to Run
package toulmin

import "testing"

func TestLevelEvaluateTracePath(t *testing.T) {
	g := NewGraph("level-eval")
	g.Counter(lvlCounter)
	results, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	var found bool
	for _, r := range results {
		for _, e := range r.Trace {
			if shortName(e.Name) != "lvlCounter" {
				continue
			}
			found = true
			if e.Level != Fail {
				t.Errorf("lvlCounter Level (Evaluate path) = %v, want Fail", e.Level)
			}
		}
	}
	if !found {
		t.Fatal("lvlCounter not found in Evaluate trace")
	}
}
