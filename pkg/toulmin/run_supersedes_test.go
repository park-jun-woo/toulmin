//ff:func feature=engine type=engine control=sequence
//ff:what TestRunSupersedes — full-pass Run marks superseded targets Excluded without polluting the Attacks graph
package toulmin

import "testing"

func TestRunSupersedes(t *testing.T) {
	ctx := NewContext()

	// Golden: X (RebuttalB, Active) supersedes Y (WarrantA). Y.Excluded == true
	// regardless of Y's own Activated/Verdict; the superseder itself is not Excluded.
	t.Run("active superseder excludes target", func(t *testing.T) {
		g := NewGraph("sup")
		y := g.Rule(WarrantA)
		x := g.Counter(RebuttalB)
		x.Supersedes(y)
		_, tr, err := g.Run(ctx)
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		ye, ok := tr.Get("WarrantA")
		if !ok {
			t.Fatal("WarrantA missing from trace")
		}
		if !ye.Excluded {
			t.Fatalf("expected WarrantA.Excluded==true, got %+v", ye)
		}
		if !ye.Activated {
			t.Fatalf("Excluded must not change Activated; expected still active, got %+v", ye)
		}
		xe, _ := tr.Get("RebuttalB")
		if xe.Excluded {
			t.Fatalf("superseder itself must not be Excluded, got %+v", xe)
		}
	})

	// An inactive superseder does not exclude its target.
	t.Run("inactive superseder does not exclude", func(t *testing.T) {
		g := NewGraph("sup2")
		y := g.Rule(WarrantA)
		x := g.Counter(InactiveR)
		x.Supersedes(y)
		_, tr, err := g.Run(ctx)
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		ye, _ := tr.Get("WarrantA")
		if ye.Excluded {
			t.Fatalf("inactive superseder must not exclude target, got %+v", ye)
		}
	})

	// Excluded does not remove Y from the Attacks graph: Y attacks Z, and Z's verdict
	// is computed with Y's attack intact even though Y is superseded (X supersedes Y).
	t.Run("supersede does not pollute Attacks graph", func(t *testing.T) {
		g := NewGraph("sup3")
		z := g.Rule(WarrantA)
		y := g.Counter(RebuttalB)
		y.Attacks(z)
		x := g.Except(DefeaterC)
		x.Supersedes(y)
		_, tr, err := g.Run(ctx)
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
		ye, _ := tr.Get("RebuttalB")
		if !ye.Excluded {
			t.Fatalf("expected RebuttalB.Excluded==true, got %+v", ye)
		}
		ze, _ := tr.Get("WarrantA")
		// Y still attacks Z: h-Categoriser yields verdict 0 (balanced), not +1.
		if ze.Verdict != 0 {
			t.Fatalf("Y's attack must still count against Z; expected verdict 0, got %v", ze.Verdict)
		}
	})
}
