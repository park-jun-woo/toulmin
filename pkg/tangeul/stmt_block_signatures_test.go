//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what stmtBlockSignatures — the ordered (Kind, Sec, Path) signatures of a Doc's non-prose blocks
package tangeul

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// stmtBlockSignatures reduces a Block tiling to the ordered structural signature the
// Phase 008 step-2 cross-check compares: each Section/Stmt block rendered as
// "Kind/Sec/Path", with ProseBlock tiles dropped. Prose is the one axis on which the
// hand-assembled step-1 golden (4-byte placeholder spans, no separate prose) and the
// tangl-en encoder (real prose spans interleaved) legitimately differ, so it is
// excluded here; what must agree is the sequence of statement/section anchors and
// their AST coordinates (StmtRef). See TestEncoderGoldenXCheck for the contract.
func stmtBlockSignatures(blocks []doc.Block) []string {
	sigs := make([]string, 0, len(blocks))
	for _, b := range blocks {
		if b.Kind == doc.ProseBlock {
			continue
		}
		sigs = append(sigs, fmt.Sprintf("%d/%d/%v", b.Kind, b.Ref.Sec, b.Ref.Path))
	}
	return sigs
}
