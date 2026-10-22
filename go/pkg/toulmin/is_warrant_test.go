//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestIsWarrant — tests isWarrant for Defeater-false, attacker-false, and non-attacker-true branches
package toulmin

import "testing"

func TestIsWarrant(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Defeater", func(t *testing.T) {
			attackerSet := map[string]bool{}
			if isWarrant(attackerSet, Defeater, "a") {
				t.Fatal("expected false for Defeater strength")
			}
		}},
		{"Attacker", func(t *testing.T) {
			attackerSet := map[string]bool{"a": true}
			if isWarrant(attackerSet, Defeasible, "a") {
				t.Fatal("expected false for attacker node")
			}
		}},
		{"NonAttacker", func(t *testing.T) {
			attackerSet := map[string]bool{}
			if !isWarrant(attackerSet, Defeasible, "a") {
				t.Fatal("expected true for non-attacker node")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
