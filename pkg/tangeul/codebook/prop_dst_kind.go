//ff:func feature=tangl type=util control=selection
//ff:what PropDstKind — the dst slot kind a PropCode implies (§2.3), do/undo→FuncRef
package codebook

import "fmt"

// PropDstKind reports the DstKind that a PropCode's Edge dst occupies (§2.3).
// The PropCode alone decides the slot layout: do/undo carry a FuncRef,
// label/importedFrom carry Text, provenance carries a Block, and every other
// predicate carries a single-word TID. An unknown PropCode is a hard error —
// never a silent fallback.
func PropDstKind(p PropCode) (DstKind, error) {
	switch p {
	case PropDo, PropUndo:
		return DstFuncRef, nil
	case PropLabel, PropImportedFrom:
		return DstText, nil
	case PropProvenance:
		return DstBlock, nil
	case PropAttacks, PropRun, PropChecks, PropRequires, PropMeans, PropHas,
		PropRuns, PropEvaluates, PropTriggers, PropUntil, PropSameAs:
		return DstTID, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown PropCode 0x%04X", uint16(p))
	}
}
