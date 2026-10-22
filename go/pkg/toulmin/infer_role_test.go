//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestInferRole — tests inferRole for except, counter, and rule branches
package toulmin

import "testing"

func TestInferRole(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Except", func(t *testing.T) {
			strMap := map[string]Strength{"a": Defeater}
			attackerSet := map[string]bool{"a": true}
			if got := inferRole(strMap, attackerSet, "a"); got != "except" {
				t.Fatalf("expected %q, got %q", "except", got)
			}
		}},
		{"Counter", func(t *testing.T) {
			strMap := map[string]Strength{"a": Defeasible}
			attackerSet := map[string]bool{"a": true}
			if got := inferRole(strMap, attackerSet, "a"); got != "counter" {
				t.Fatalf("expected %q, got %q", "counter", got)
			}
		}},
		{"Rule", func(t *testing.T) {
			strMap := map[string]Strength{"a": Strict}
			attackerSet := map[string]bool{}
			if got := inferRole(strMap, attackerSet, "a"); got != "rule" {
				t.Fatalf("expected %q, got %q", "rule", got)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
