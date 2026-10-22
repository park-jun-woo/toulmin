//ff:func feature=tangl type=util control=selection
//ff:what sectionHeader — the `## tangl:<name>` heading line for a locale and SectionCode
package doc

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// sectionHeader returns the mode-2 section heading line for the given locale and
// SectionCode. The markdown skeleton is the core renderer's concern (README
// contract §3); the heading keyword is locale specific (EN `Cases` <-> KO `사안`,
// files/TANGL문법v0.3-한국어.md §일곱 섹션). An unknown locale or section is an error.
func sectionHeader(locale string, sec toulmin.SectionCode) (string, error) {
	var names map[toulmin.SectionCode]string
	switch locale {
	case "en":
		names = map[toulmin.SectionCode]string{
			toulmin.SectionSubject:     "Subject",
			toulmin.SectionSee:         "See",
			toulmin.SectionDefinitions: "Definitions",
			toulmin.SectionRules:       "Rules",
			toulmin.SectionCases:       "Cases",
			toulmin.SectionProvides:    "Provides",
			toulmin.SectionInternal:    "Internal",
		}
	case "ko":
		names = map[toulmin.SectionCode]string{
			toulmin.SectionSubject:     "주제",
			toulmin.SectionSee:         "참조",
			toulmin.SectionDefinitions: "정의",
			toulmin.SectionRules:       "규칙",
			toulmin.SectionCases:       "사안",
			toulmin.SectionProvides:    "업무",
			toulmin.SectionInternal:    "내부",
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
