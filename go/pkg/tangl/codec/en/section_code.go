//ff:func feature=tangl type=parser control=selection
//ff:what sectionCode — map a verbatim `## tangl:<Name>` heading name to its codebook SectionCode
package en

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// sectionCode maps a scanner Section.Name (the verbatim `## tangl:<Name>` label)
// to its codebook SectionCode. The seven names match the parser's recognized set;
// ok is false for any other name (a scanner/parser drift the caller reports).
func sectionCode(name string) (toulmin.SectionCode, bool) {
	switch name {
	case "Subject":
		return toulmin.SectionSubject, true
	case "See":
		return toulmin.SectionSee, true
	case "Definitions":
		return toulmin.SectionDefinitions, true
	case "Rules":
		return toulmin.SectionRules, true
	case "Cases":
		return toulmin.SectionCases, true
	case "Provides":
		return toulmin.SectionProvides, true
	case "Internal":
		return toulmin.SectionInternal, true
	default:
		return 0, false
	}
}
