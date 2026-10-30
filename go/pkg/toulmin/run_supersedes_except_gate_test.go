//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunSupersedesExceptGate — a superseder excludes only while a_eff > 0 (Except gates, Counter does not)
package toulmin

import "testing"

// TestRunSupersedesExceptGate is the cross-port VIP scenario: vipDiscount supersedes
// generalDiscount. An exception (suspendedVIP) that fully gates vipDiscount (a_eff 0)
// lifts the exclusion; a partial exception (F 0.5, a_eff 0.5) or a rebuttal (vipDisputed,
// a_eff stays 1) keeps it. Both Run and Evaluate(Trace) must agree.
func TestRunSupersedesExceptGate(t *testing.T) {
	cases := []struct {
		name         string
		wire         func(g *Graph, vip *Rule)
		vipVerdict   float64
		vipAEff      float64
		wantExcluded bool
	}{
		{"no exception", func(g *Graph, vip *Rule) {}, 1, 1, true},
		{"exception fully gates", func(g *Graph, vip *Rule) { g.Except(suspendedVIP).Attacks(vip) }, -1, 0, false},
		{"exception force 0.5", func(g *Graph, vip *Rule) { g.Except(suspendedVIP).Attacks(vip, 0.5) }, 0, 0.5, true},
		{"counter keeps priority", func(g *Graph, vip *Rule) { g.Counter(vipDisputed).Attacks(vip) }, 0, 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			checkVIPSupersede(t, tc.wire, tc.vipVerdict, tc.vipAEff, tc.wantExcluded)
		})
	}
}
