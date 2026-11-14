//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRuleLevel — tests Rule.Level sets Level+LevelSet and chains; Unset stays explicit
package toulmin

import "testing"

func TestRuleLevel(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"SetsFieldsAndChains", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			got := r.Level(Review)
			if got != r {
				t.Errorf("Level must return the receiver for chaining, got %v want %v", got, r)
			}
			if g.rules[r.idx].Level != Review {
				t.Errorf("expected Level=Review, got %v", g.rules[r.idx].Level)
			}
			if !g.rules[r.idx].LevelSet {
				t.Error("expected LevelSet=true after Level call")
			}
		}},
		{"ExplicitUnsetIsSet", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			r.Level(Unset)
			if g.rules[r.idx].Level != Unset {
				t.Errorf("expected Level=Unset, got %v", g.rules[r.idx].Level)
			}
			if !g.rules[r.idx].LevelSet {
				t.Error("explicit Level(Unset) must still mark LevelSet=true")
			}
		}},
		{"NotSetByDefault", func(t *testing.T) {
			g := NewGraph("test")
			r := g.Counter(RebuttalB)
			if g.rules[r.idx].LevelSet {
				t.Error("LevelSet must be false without a Level call")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
