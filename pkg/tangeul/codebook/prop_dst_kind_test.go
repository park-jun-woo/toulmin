//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestPropDstKind — do/undo fixed to FuncRef, Text/Block/TID per §2.3, unknown errors
package codebook

import "testing"

func TestPropDstKind(t *testing.T) {
	cases := []struct {
		prop PropCode
		want DstKind
	}{
		// do/undo are the only FuncRef dst — the invariant this test locks.
		{PropDo, DstFuncRef},
		{PropUndo, DstFuncRef},
		// label/importedFrom carry Text.
		{PropLabel, DstText},
		{PropImportedFrom, DstText},
		// provenance carries a Block.
		{PropProvenance, DstBlock},
		// everything else is a single-word TID.
		{PropAttacks, DstTID},
		{PropRun, DstTID},
		{PropChecks, DstTID},
		{PropRequires, DstTID},
		{PropMeans, DstTID},
		{PropHas, DstTID},
		{PropRuns, DstTID},
		{PropEvaluates, DstTID},
		{PropTriggers, DstTID},
		{PropUntil, DstTID},
		{PropSameAs, DstTID},
		{PropSupersedes, DstTID},
	}
	for _, c := range cases {
		got, err := PropDstKind(c.prop)
		if err != nil || got != c.want {
			t.Errorf("PropDstKind(0x%04X): got %d err=%v, want %d", uint16(c.prop), got, err, c.want)
		}
	}
	if _, err := PropDstKind(PropCode(0xFFFF)); err == nil {
		t.Error("PropDstKind(unknown): expected error")
	}
}
