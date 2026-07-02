//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestNewEvalContext — tests newEvalContext for nil/non-nil defeatEdges and cycle error branches
package toulmin

import "testing"

func TestNewEvalContext(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NilDefeatEdges", func(t *testing.T) {
			rules := []RuleMeta{
				{Name: "WarrantA", Qualifier: 1.0, Strength: Strict, Fn: WarrantA},
			}
			roleMap := map[string]string{"WarrantA": "warrant"}
			ec, err := newEvalContext(rules, nil, roleMap)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ec.fnMap["WarrantA"] == nil {
				t.Fatalf("expected fnMap populated")
			}
			if ec.qualMap["WarrantA"] != 1.0 {
				t.Fatalf("expected qualMap populated")
			}
			if ec.strMap["WarrantA"] != Strict {
				t.Fatalf("expected strMap populated")
			}
		}},
		{"WithDefeatEdges", func(t *testing.T) {
			rules := []RuleMeta{
				{Name: "WarrantA", Qualifier: 1.0, Strength: Strict, Fn: WarrantA},
				{Name: "RebuttalB", Qualifier: 0.5, Strength: Defeasible, Fn: RebuttalB},
			}
			edges := []defeatEdge{{from: "RebuttalB", to: "WarrantA"}}
			ec, err := newEvalContext(rules, edges, map[string]string{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(ec.edges["WarrantA"]) != 1 || ec.edges["WarrantA"][0] != "RebuttalB" {
				t.Fatalf("expected edges built from defeatEdges, got %v", ec.edges)
			}
			if !ec.attackerSet["RebuttalB"] {
				t.Fatalf("expected RebuttalB in attackerSet")
			}
		}},
		{"CycleError", func(t *testing.T) {
			rules := []RuleMeta{
				{Name: "A", Qualifier: 1.0, Strength: Strict, Fn: WarrantA},
				{Name: "B", Qualifier: 1.0, Strength: Strict, Fn: WarrantA},
			}
			edges := []defeatEdge{
				{from: "A", to: "B"},
				{from: "B", to: "A"},
			}
			ec, err := newEvalContext(rules, edges, map[string]string{})
			if err == nil {
				t.Fatalf("expected cycle error")
			}
			if ec != nil {
				t.Fatalf("expected nil evalContext on error")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
