//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what matchBodyRefs — index the Line-less Provides/Internal run·check bodies into refByLine
package en

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// matchBodyRefs fills refByLine for the Line-less run/check bodies of Provides and
// Internal (their Runs/Checks are `[]string`, README decision 3). It dispatches each
// such section to matchSectionBodies, which walks the scanned body items and assigns
// per-kind slot indices; non-body sections are skipped.
func matchBodyRefs(a *ast.Document, secs []scan.Section, refByLine map[int]doc.StmtRef) error {
	for _, sec := range secs {
		code, ok := sectionCode(sec.Name)
		if !ok || (code != codebook.SectionProvides && code != codebook.SectionInternal) {
			continue
		}
		if err := matchSectionBodies(a, code, sec, refByLine); err != nil {
			return err
		}
	}
	return nil
}
