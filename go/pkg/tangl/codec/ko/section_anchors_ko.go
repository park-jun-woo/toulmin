//ff:func feature=tangl type=parser control=selection
//ff:what sectionAnchorsKo — dispatch a section's items to its per-section anchor walker
package ko

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// sectionAnchorsKo dispatches one section's scanned items to the per-section walker
// for its codebook SectionCode, each of which parses the items into the Document and
// returns the section's StmtBlock anchors.
func sectionAnchorsKo(a *ast.Document, code toulmin.SectionCode, items []scan.Item, path string) ([]doc.Block, error) {
	switch code {
	case toulmin.SectionSubject:
		return subjectAnchorsKo(a, items, path)
	case toulmin.SectionSee:
		return seeAnchorsKo(a, items, path)
	case toulmin.SectionDefinitions:
		return definitionAnchorsKo(a, items, path)
	case toulmin.SectionRules:
		return ruleAnchorsKo(a, items, path)
	case toulmin.SectionCases:
		return caseAnchorsKo(a, items, path)
	case toulmin.SectionProvides:
		return providesAnchorsKo(a, items, path)
	default:
		return internalAnchorsKo(a, items, path)
	}
}
