//ff:func feature=tangl type=parser control=selection
//ff:what sectionCodeKo — map a canonical English section name to its codebook SectionCode
package ko

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// sectionCodeKo maps the canonical English section name (as returned by
// koSectionNames from a Korean `## tangl:<이름>` heading) to its codebook
// SectionCode. Because the codes are 1..7 in required document order, Encode also
// uses them to enforce section ordering. ok is false for any other name.
func sectionCodeKo(canon string) (toulmin.SectionCode, bool) {
	switch canon {
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
