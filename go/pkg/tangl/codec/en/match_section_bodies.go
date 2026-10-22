//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what matchSectionBodies — match one Provides/Internal section's endpoints to their body slots
package en

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// matchSectionBodies walks the scanned top-level items of one Provides or Internal
// section, resolves each to its AST endpoint/internal by source line, and delegates
// the run/check body assignment to assignBodySlots. A non-struck head with no AST
// match is a scanner/parser drift reported as error.
func matchSectionBodies(a *ast.Document, code toulmin.SectionCode, sec scan.Section, refByLine map[int]doc.StmtRef) error {
	items, _ := scan.ScanItems(sec.Lines)
	for _, head := range items {
		if head.Struck {
			continue
		}
		runs, checks, idx, ok := bodySlices(a, code, head.No)
		if !ok {
			return fmt.Errorf("tangl/codec/en: no AST %s entry for item at line %d", sec.Name, head.No)
		}
		if err := assignBodySlots(head.Children, code, idx, runs, checks, refByLine); err != nil {
			return err
		}
	}
	return nil
}
