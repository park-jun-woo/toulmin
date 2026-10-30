//ff:func feature=engine type=engine control=sequence
//ff:what checkVIPSupersede — test helper: builds the VIP supersede graph, checks Run and Evaluate(Trace) agree
package toulmin

import "testing"

// checkVIPSupersede builds generalDiscount ← Supersedes ← vipDiscount, applies wire, and
// checks generalDiscount's Excluded flag (Run and Evaluate Trace), vipDiscount's verdict,
// and vipDiscount's a_eff.
func checkVIPSupersede(t *testing.T, wire func(g *Graph, vip *Rule), vipVerdict, vipAEff float64, wantExcluded bool) {
	t.Helper()
	g := NewGraph("vip")
	general := g.Rule(generalDiscount)
	vip := g.Rule(vipDiscount)
	vip.Supersedes(general)
	wire(g, vip)

	results, tr, err := g.Run(NewContext())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	ge, _ := tr.Get("generalDiscount")
	ve, _ := tr.Get("vipDiscount")
	if ge.Excluded != wantExcluded {
		t.Fatalf("generalDiscount.Excluded = %v, want %v", ge.Excluded, wantExcluded)
	}
	if ve.Excluded || !ve.Activated || ve.Verdict != vipVerdict {
		t.Fatalf("vipDiscount entry = %+v, want not excluded, activated, verdict %v", ve, vipVerdict)
	}
	if ge.Verdict != 1 {
		t.Fatalf("Excluded is reporting-only; generalDiscount verdict must stay 1, got %v", ge.Verdict)
	}
	if got := applicabilityOf(results, "vipDiscount"); got != vipAEff {
		t.Fatalf("vipDiscount a_eff = %v, want %v", got, vipAEff)
	}

	traced, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if ex, ok := excludedIn(traced, "generalDiscount", "generalDiscount"); !ok || ex != wantExcluded {
		t.Fatalf("Evaluate trace: generalDiscount Excluded = %v (ok=%v), want %v", ex, ok, wantExcluded)
	}
}
