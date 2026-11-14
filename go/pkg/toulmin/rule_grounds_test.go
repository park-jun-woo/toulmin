//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRuleGrounds — tests Rule.Grounds records/appends names and leaves Specs/Level/Defeats/ruleID untouched
package toulmin

import "testing"

func TestRuleGrounds(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"SingleCallRecordsAndChains", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			got := r.Grounds("a")
			if got != r {
				t.Errorf("Grounds must return the receiver for chaining, got %v want %v", got, r)
			}
			grounds := g.rules[r.idx].Grounds
			if len(grounds) != 1 || grounds[0] != "a" {
				t.Errorf("expected Grounds=[a], got %v", grounds)
			}
		}},
		{"MultipleCallsAppend", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			r.Grounds("a").Grounds("b")
			grounds := g.rules[r.idx].Grounds
			if len(grounds) != 2 || grounds[0] != "a" || grounds[1] != "b" {
				t.Errorf("expected Grounds=[a b], got %v", grounds)
			}
		}},
		{"VariadicAppends", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			r.Grounds("source-body", "mx")
			grounds := g.rules[r.idx].Grounds
			if len(grounds) != 2 || grounds[0] != "source-body" || grounds[1] != "mx" {
				t.Errorf("expected Grounds=[source-body mx], got %v", grounds)
			}
		}},
		{"DoesNotDisturbOtherFields", func(t *testing.T) {
			g := NewGraph("test")
			w := g.Rule(WarrantA)
			r := g.Counter(RebuttalB)
			r.Attacks(w)
			r.Level(Review)

			oldID := r.id
			oldName := g.rules[r.idx].Name
			oldSpecsLen := len(g.rules[r.idx].Specs)
			oldLevel := g.rules[r.idx].Level
			oldDefeats := len(g.defeats)

			r.Grounds("a").Grounds("b")

			if g.rules[r.idx].Name != oldName {
				t.Errorf("Grounds must not change Name: was %q now %q", oldName, g.rules[r.idx].Name)
			}
			if r.id != oldID {
				t.Errorf("Grounds must not change ruleID: was %q now %q", oldID, r.id)
			}
			if len(g.rules[r.idx].Specs) != oldSpecsLen {
				t.Errorf("Grounds must not change Specs, len was %d now %d", oldSpecsLen, len(g.rules[r.idx].Specs))
			}
			if g.rules[r.idx].Level != oldLevel {
				t.Errorf("Grounds must not change Level: was %v now %v", oldLevel, g.rules[r.idx].Level)
			}
			if len(g.defeats) != oldDefeats {
				t.Errorf("Grounds must not change defeats, len was %d now %d", oldDefeats, len(g.defeats))
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
