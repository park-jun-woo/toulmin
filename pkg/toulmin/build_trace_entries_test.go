//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestBuildTraceEntries — tests Graph.buildTraceEntries for no-rules and multi-rule branches
package toulmin

import "testing"

func TestBuildTraceEntries(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NoRules", func(t *testing.T) {
			g := &Graph{name: "g", roles: map[string]string{}}
			ec := &evalContext{
				strMap:       map[string]Strength{},
				attackerSet:  map[string]bool{},
				active:       map[string]bool{},
				qualMap:      map[string]float64{},
				verdictCache: map[string]float64{},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
			}
			entries := g.buildTraceEntries(ec, NewContext())
			if len(entries) != 0 {
				t.Fatalf("expected 0 entries, got %d", len(entries))
			}
		}},
		{"MultiRules", func(t *testing.T) {
			g := &Graph{
				name:  "g",
				roles: map[string]string{"WarrantA": "rule"},
				rules: []RuleMeta{
					{Name: "WarrantA", Qualifier: 1.0, Strength: Strict, Fn: WarrantA},
					{Name: "RebuttalB", Qualifier: 0.8, Strength: Defeasible, Fn: RebuttalB},
				},
			}
			ec := &evalContext{
				strMap:       map[string]Strength{"WarrantA": Strict, "RebuttalB": Defeasible},
				attackerSet:  map[string]bool{"RebuttalB": true},
				active:       map[string]bool{"WarrantA": true, "RebuttalB": true},
				qualMap:      map[string]float64{"WarrantA": 1.0, "RebuttalB": 0.8},
				verdictCache: map[string]float64{"WarrantA": 1.0, "RebuttalB": 0.8},
				evidence:     map[string]any{},
				specsMap:     map[string]Specs{},
			}
			entries := g.buildTraceEntries(ec, NewContext())
			if len(entries) != 2 {
				t.Fatalf("expected 2 entries, got %d", len(entries))
			}
			if entries[0].Name != "WarrantA" || entries[0].Role != "rule" {
				t.Errorf("entries[0] unexpected: %+v", entries[0])
			}
			if entries[1].Name != "RebuttalB" || entries[1].Role == "" {
				t.Errorf("entries[1] unexpected: %+v", entries[1])
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
