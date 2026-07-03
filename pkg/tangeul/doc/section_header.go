//ff:func feature=tangl type=util control=selection
//ff:what sectionHeader — the `## tangl:<name>` heading line for a locale and SectionCode
package doc

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

// sectionHeader returns the mode-2 section heading line for the given locale and
// SectionCode. The markdown skeleton is the core renderer's concern (README
// contract §3); the heading keyword is locale specific (EN `Cases` <-> KO `사안`,
// files/TANGL문법v0.3-한국어.md §일곱 섹션). An unknown locale or section is an error.
func sectionHeader(locale string, sec codebook.SectionCode) (string, error) {
	var names map[codebook.SectionCode]string
	switch locale {
	case "en":
		names = map[codebook.SectionCode]string{
			codebook.SectionSubject:     "Subject",
			codebook.SectionSee:         "See",
			codebook.SectionDefinitions: "Definitions",
			codebook.SectionRules:       "Rules",
			codebook.SectionCases:       "Cases",
			codebook.SectionProvides:    "Provides",
			codebook.SectionInternal:    "Internal",
		}
	case "ko":
		names = map[codebook.SectionCode]string{
			codebook.SectionSubject:     "주제",
			codebook.SectionSee:         "참조",
			codebook.SectionDefinitions: "정의",
			codebook.SectionRules:       "규칙",
			codebook.SectionCases:       "사안",
			codebook.SectionProvides:    "업무",
			codebook.SectionInternal:    "내부",
		}
	default:
		return "", fmt.Errorf("tangeul/doc: no section headers for locale %q", locale)
	}
	name, ok := names[sec]
	if !ok {
		return "", fmt.Errorf("tangeul/doc: unknown SectionCode %d for locale %q", uint16(sec), locale)
	}
	return "## tangl:" + name, nil
}
