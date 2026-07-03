//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestSectionCodeKo — all seven canonical names map to their codebook SectionCode; unknown reports false
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
)

func TestSectionCodeKo(t *testing.T) {
	want := map[string]codebook.SectionCode{
		"Subject":     codebook.SectionSubject,
		"See":         codebook.SectionSee,
		"Definitions": codebook.SectionDefinitions,
		"Rules":       codebook.SectionRules,
		"Cases":       codebook.SectionCases,
		"Provides":    codebook.SectionProvides,
		"Internal":    codebook.SectionInternal,
	}
	for canon, code := range want {
		got, ok := sectionCodeKo(canon)
		if !ok || got != code {
			t.Errorf("sectionCodeKo(%q) = (%v, %v), want (%v, true)", canon, got, ok, code)
		}
	}

	if got, ok := sectionCodeKo("주제"); ok || got != 0 {
		t.Fatalf("sectionCodeKo(%q) = (%v, %v), want (0, false)", "주제", got, ok)
	}
}
