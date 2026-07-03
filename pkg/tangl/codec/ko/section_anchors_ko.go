//ff:func feature=tangl type=parser control=selection
//ff:what sectionAnchorsKo — dispatch a section's items to its per-section anchor walker
package ko

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// sectionAnchorsKo dispatches one section's scanned items to the per-section walker
// for its codebook SectionCode, each of which parses the items into the Document and
// returns the section's StmtBlock anchors.
func sectionAnchorsKo(a *ast.Document, code codebook.SectionCode, items []scan.Item, path string) ([]doc.Block, error) {
	switch code {
	case codebook.SectionSubject:
		return subjectAnchorsKo(a, items, path)
	case codebook.SectionSee:
		return seeAnchorsKo(a, items, path)
	case codebook.SectionDefinitions:
		return definitionAnchorsKo(a, items, path)
	case codebook.SectionRules:
		return ruleAnchorsKo(a, items, path)
	case codebook.SectionCases:
		return caseAnchorsKo(a, items, path)
	case codebook.SectionProvides:
		return providesAnchorsKo(a, items, path)
	default:
		return internalAnchorsKo(a, items, path)
	}
}
