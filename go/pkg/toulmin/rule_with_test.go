//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRuleWith — tests Rule.With for validation panic, role rename, missing-role skip, and defeats-edge rewrite branches
package toulmin

import (
	"strings"
	"testing"
)

func TestRuleWith(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"PanicsOnInvalidSpec", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic for invalid spec")
				}
				msg, ok := r.(string)
				if !ok {
					t.Fatalf("expected string panic, got %T", r)
				}
				if !strings.Contains(msg, "toulmin:") {
					t.Errorf("unexpected panic message: %s", msg)
				}
			}()
			g := NewGraph("test")
			g.Rule(WarrantA).With(&ruleIDUnmarshalableSpec{})
		}},
		{"RenamesRoleAndDefeatsEdges", func(t *testing.T) {
			g := NewGraph("test")
			w := g.Rule(WarrantA)
			r := g.Counter(RebuttalB)
			r.Attacks(w) // defeats: from=r.id(oldID) to=w.id
			w.Attacks(r) // defeats: from=w.id to=r.id(oldID) -> exercises "to" match branch too

			oldID := r.id
			if _, ok := g.roles[oldID]; !ok {
				t.Fatalf("expected role registered for %q before With", oldID)
			}

			r.With(&testSpec{Value: "x"})
			newID := r.id
			if newID == oldID {
				t.Fatalf("expected id to change after With")
			}

			if _, ok := g.roles[oldID]; ok {
				t.Errorf("expected old role entry removed for %q", oldID)
			}
			if role, ok := g.roles[newID]; !ok || role != "counter" {
				t.Errorf("expected role migrated to newID %q with role counter, got %q ok=%v", newID, role, ok)
			}

			foundFrom, foundTo := false, false
			for _, d := range g.defeats {
				if d.from == newID {
					foundFrom = true
				}
				if d.to == newID {
					foundTo = true
				}
				if d.from == oldID || d.to == oldID {
					t.Errorf("expected no defeat edges referencing stale oldID %q, got %+v", oldID, d)
				}
			}
			if !foundFrom {
				t.Errorf("expected a defeat edge with from=newID")
			}
			if !foundTo {
				t.Errorf("expected a defeat edge with to=newID")
			}
		}},
		{"SkipsWhenRoleMissing", func(t *testing.T) {
			fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			g := NewGraph("test")
			wrapped := toRuleFunc(fn)
			id := ruleID(fn, nil)
			g.rules = append(g.rules, RuleMeta{
				Name:      id,
				Qualifier: 1.0,
				Strength:  Defeasible,
				Fn:        wrapped,
			})
			// Intentionally do not register g.roles[id], to exercise the "ok" == false branch.
			r := &Rule{id: id, graph: g, idx: 0, fn: fn}

			got := r.With(&testSpec{Value: "z"})
			if got != r {
				t.Errorf("With must return the receiver for chaining, got %v want %v", got, r)
			}
			if len(g.roles) != 0 {
				t.Errorf("expected roles map to remain empty, got %v", g.roles)
			}
			if g.rules[0].Name != r.id {
				t.Errorf("expected rules[0].Name updated to new id %q, got %q", r.id, g.rules[0].Name)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
