//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestBuildTraceEntry — tests Graph.buildTraceEntry for explicit-role and inferred-role branches
package toulmin

import "testing"

func TestBuildTraceEntry(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"ExplicitRole", func(t *testing.T) {
			g := &Graph{roles: map[string]string{"WarrantA": "rule"}}
			ec := &evalContext{
				strMap:       map[string]Strength{"WarrantA": Strict},
				attackerSet:  map[string]bool{},
				active:       map[string]bool{"WarrantA": true},
				qualMap:      map[string]float64{"WarrantA": 1.0},
				verdictCache: map[string]float64{"WarrantA": 1.0},
				evidence:     map[string]any{"WarrantA": "ev"},
				specsMap:     map[string]Specs{"WarrantA": nil},
			}
			entry := g.buildTraceEntry(ec, "WarrantA", NewContext())
			if entry.Name != "WarrantA" {
				t.Errorf("Name = %q, want %q", entry.Name, "WarrantA")
			}
			if entry.Role != "rule" {
				t.Errorf("Role = %q, want %q", entry.Role, "rule")
			}
			if !entry.Activated {
				t.Error("expected Activated = true")
			}
			if entry.Qualifier != 1.0 {
				t.Errorf("Qualifier = %v, want 1.0", entry.Qualifier)
			}
			if entry.Verdict != 1.0 {
				t.Errorf("Verdict = %v, want 1.0", entry.Verdict)
			}
			if entry.Evidence != "ev" {
				t.Errorf("Evidence = %v, want %q", entry.Evidence, "ev")
			}
		}},
		{"InferredRole", func(t *testing.T) {
			g := &Graph{roles: map[string]string{}}
			ec := &evalContext{
				strMap:       map[string]Strength{"RebuttalB": Defeasible},
				attackerSet:  map[string]bool{"RebuttalB": true},
				active:       map[string]bool{"RebuttalB": false},
				qualMap:      map[string]float64{"RebuttalB": 0.5},
				verdictCache: map[string]float64{"RebuttalB": 0.5},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
			}
			entry := g.buildTraceEntry(ec, "RebuttalB", NewContext())
			if entry.Name != "RebuttalB" {
				t.Errorf("Name = %q, want %q", entry.Name, "RebuttalB")
			}
			if entry.Role == "" {
				t.Error("expected inferred non-empty role")
			}
			if entry.Activated {
				t.Error("expected Activated = false")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
