//ff:func feature=engine type=engine control=sequence
//ff:what TestRunSupersedesChain — X>Y>Z chain: an excluded superseder still excludes; a gated one does not
package toulmin

import "testing"

// TestRunSupersedesChain: X supersedes Y, Y supersedes Z. Exclusion is reporting-only, so
// an excluded Y keeps a_eff 1 and still excludes Z. When an exception fully gates X
// (a_eff 0), only X loses its power: Y is no longer excluded, Z still is.
func TestRunSupersedesChain(t *testing.T) {
	build := func(gateX bool) Trace {
		g := NewGraph("chain")
		z := g.Rule(DefeaterC)
		y := g.Rule(RebuttalB)
		x := g.Rule(WarrantA)
		x.Supersedes(y)
		y.Supersedes(z)
		if gateX {
			g.Except(suspendedVIP).Attacks(x)
		}
		_, tr, err := g.Run(NewContext())
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		return tr
	}
	excl := func(tr Trace, name string) bool {
		e, _ := tr.Get(name)
		return e.Excluded
	}

	tr := build(false)
	if excl(tr, "WarrantA") || !excl(tr, "RebuttalB") || !excl(tr, "DefeaterC") {
		t.Fatalf("ungated chain: want X kept, Y and Z excluded, got X=%v Y=%v Z=%v",
			excl(tr, "WarrantA"), excl(tr, "RebuttalB"), excl(tr, "DefeaterC"))
	}
	tr = build(true)
	if excl(tr, "WarrantA") || excl(tr, "RebuttalB") || !excl(tr, "DefeaterC") {
		t.Fatalf("gated X: want Y kept, Z excluded, got X=%v Y=%v Z=%v",
			excl(tr, "WarrantA"), excl(tr, "RebuttalB"), excl(tr, "DefeaterC"))
	}
}
