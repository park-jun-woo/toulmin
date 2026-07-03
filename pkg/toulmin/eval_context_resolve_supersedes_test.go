//ff:func feature=engine type=engine control=sequence
//ff:what TestResolveSupersedes — active superseder excludes its target; inactive one does not
package toulmin

import "testing"

func TestResolveSupersedes(t *testing.T) {
	ctx := NewContext()

	// Active superseder (RebuttalB) supersedes WarrantA → WarrantA excluded.
	g := NewGraph("resolve")
	w := g.Rule(WarrantA)
	s := g.Counter(RebuttalB)
	s.Supersedes(w)
	_, ec, err := g.evaluate(ctx, EvalOption{}, true)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	excluded := ec.resolveSupersedes(ctx)
	if !excluded[w.id] {
		t.Fatalf("expected WarrantA excluded by active superseder, got %v", excluded)
	}

	// Inactive superseder (InactiveR) supersedes WarrantA → WarrantA not excluded.
	g2 := NewGraph("resolve2")
	w2 := g2.Rule(WarrantA)
	s2 := g2.Counter(InactiveR)
	s2.Supersedes(w2)
	_, ec2, err := g2.evaluate(ctx, EvalOption{}, true)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	excluded2 := ec2.resolveSupersedes(ctx)
	if excluded2[w2.id] {
		t.Fatalf("inactive superseder must not exclude target, got %v", excluded2)
	}

	// resolveSupersedes must force an unreached superseder to run (idempotent calc).
	// InactiveR here is never a warrant/attacker, but a forced calc still marks it ran.
	if !ec2.ran[s2.id] {
		t.Fatalf("resolveSupersedes must force-run superseder %q", s2.id)
	}
}
