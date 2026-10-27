//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestSectionCode — map every section name and reject an unknown one
package en

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// TestSectionCode checks the seven recognized section names map to their codebook
// SectionCode in declaration order and that an unknown name reports ok=false.
func TestSectionCode(t *testing.T) {
	want := []struct {
		name string
		code toulmin.SectionCode
	}{
		{"Subject", toulmin.SectionSubject},
		{"See", toulmin.SectionSee},
		{"Definitions", toulmin.SectionDefinitions},
		{"Rules", toulmin.SectionRules},
		{"Cases", toulmin.SectionCases},
		{"Provides", toulmin.SectionProvides},
		{"Internal", toulmin.SectionInternal},
	}
	for _, w := range want {
		got, ok := sectionCode(w.name)
		if !ok || got != w.code {
			t.Errorf("sectionCode(%q) = (%d, %v), want (%d, true)", w.name, got, ok, w.code)
		}
	}
	if _, ok := sectionCode("Unknown"); ok {
		t.Errorf("sectionCode(Unknown) ok = true, want false")
	}
}
