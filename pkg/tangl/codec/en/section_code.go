//ff:func feature=tangl type=parser control=selection
//ff:what sectionCode — map a verbatim `## tangl:<Name>` heading name to its codebook SectionCode
package en

import "github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"

// sectionCode maps a scanner Section.Name (the verbatim `## tangl:<Name>` label)
// to its codebook SectionCode. The seven names match the parser's recognized set;
// ok is false for any other name (a scanner/parser drift the caller reports).
func sectionCode(name string) (codebook.SectionCode, bool) {
	switch name {
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
