//ff:func feature=engine type=engine control=sequence
//ff:what TestRunOnActiveOnly — RunOn fires only on Active nodes, never on Defeated or Inactive
package toulmin

import "testing"

func TestRunOnActiveOnly(t *testing.T) {
	fired := map[string]bool{}
	mark := func(name string) NodeHandler {
		return func(t Trace) error {
			fired[name] = true
			return nil
		}
	}

	// Distinct always-true rule fns (two rules in one graph need distinct func values).
	target := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	attacker := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	g := NewGraph("active-only")
	// Active: unattacked warrant → verdict 1 → Active.
	g.Rule(WarrantA).RunOn(mark("WarrantA"))
	// Defeated: target attacked by an equal-strength active counter → verdict 0 → Defeated.
	w := g.Rule(target).RunOn(mark("target"))
	g.Counter(attacker).RunOn(mark("attacker")).Attacks(w)
	// Inactive: rule fn returns false → Inactive.
	g.Rule(InactiveR).RunOn(mark("InactiveR"))

	if _, _, err := g.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !fired["WarrantA"] {
		t.Error("Active node must fire RunOn")
	}
	if !fired["attacker"] {
		t.Error("the attacker is itself an unattacked Active node and must fire RunOn")
	}
	if fired["target"] {
		t.Error("Defeated node (verdict 0) must NOT fire RunOn")
	}
	if fired["InactiveR"] {
		t.Error("Inactive node must NOT fire RunOn")
	}
}
