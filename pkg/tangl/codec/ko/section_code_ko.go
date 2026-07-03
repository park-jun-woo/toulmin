//ff:func feature=tangl type=parser control=selection
//ff:what sectionCodeKo — map a canonical English section name to its codebook SectionCode
package ko

import "github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"

// sectionCodeKo maps the canonical English section name (as returned by
// koSectionNames from a Korean `## tangl:<이름>` heading) to its codebook
// SectionCode. Because the codes are 1..7 in required document order, Encode also
// uses them to enforce section ordering. ok is false for any other name.
func sectionCodeKo(canon string) (codebook.SectionCode, bool) {
	switch canon {
	case "Subject":
		return codebook.SectionSubject, true
	case "See":
		return codebook.SectionSee, true
	case "Definitions":
		return codebook.SectionDefinitions, true
	case "Rules":
		return codebook.SectionRules, true
	case "Cases":
		return codebook.SectionCases, true
	case "Provides":
		return codebook.SectionProvides, true
	case "Internal":
		return codebook.SectionInternal, true
	default:
		return 0, false
	}
}
