//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphRunRuleVerdict — Run and Evaluate traces carry each rule's own verdict; partial rules are Active and fire handlers
package toulmin

import "testing"

func TestGraphRunRuleVerdict(t *testing.T) {
	fired := map[string]float64{}
	g := NewGraph("rv")
	w := g.Rule(func(ctx Context, specs Specs) (float64, any) { return 0.8, nil }).RunOn(func(self TraceEntry, tr Trace) error {
		fired["w"] = self.RuleVerdict
		return nil
	})
	g.Counter(func(ctx Context, specs Specs) (float64, any) { return -1.0, nil }).Attacks(w)
	_, tr, err := g.Run(NewContext())
	if err != nil {
		t.Fatal(err)
	}
	if fired["w"] != 0.8 {
		t.Fatalf("expected partial rule to fire with RuleVerdict 0.8, got %v", fired)
	}
	want := map[string]float64{"rule": 0.8, "counter": -1.0}
	for _, e := range tr.All() {
		if e.RuleVerdict != want[e.Role] || e.Activated != (e.RuleVerdict > -1) {
			t.Errorf("%s (%s): RuleVerdict %v Activated %v", e.Name, e.Role, e.RuleVerdict, e.Activated)
		}
	}
	res, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range res[0].Trace {
		if e.RuleVerdict != want[e.Role] {
			t.Errorf("evaluate trace %s (%s): RuleVerdict %v", e.Name, e.Role, e.RuleVerdict)
		}
	}
}
