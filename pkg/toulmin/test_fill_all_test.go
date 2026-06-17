//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestFillAll — fillAll runs calc only on nodes not yet evaluated
package toulmin

import "testing"

func TestFillAll(t *testing.T) {
	g := NewGraph("fill")
	g.Rule(WarrantA)
	g.Rule(InactiveR)
	ec, err := newEvalContext(g.rules, g.defeats, g.roles)
	if err != nil {
		t.Fatalf("newEvalContext: %v", err)
	}
	ctx := NewContext()

	// Pre-run one node so fillAll must skip the already-ran branch for it.
	ec.calc(g.rules[0].Name, ctx)

	ec.fillAll(g.rules, ctx)
	if ec.err != nil {
		t.Fatalf("unexpected fillAll error: %v", ec.err)
	}
	for _, r := range g.rules {
		if !ec.ran[r.Name] {
			t.Errorf("node %q was not run after fillAll", r.Name)
		}
	}
}
