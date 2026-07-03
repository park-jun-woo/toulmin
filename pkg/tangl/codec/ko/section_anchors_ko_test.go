//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestSectionAnchorsKo — each SectionCode dispatches to its per-section walker, unmapped code hits the default (Internal) walker
package ko

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

func TestSectionAnchorsKo(t *testing.T) {
	cases := []struct {
		name  string
		code  codebook.SectionCode
		items []scan.Item
		check func(t *testing.T, a *ast.Document)
	}{
		{
			name:  "Subject",
			code:  codebook.SectionSubject,
			items: []scan.Item{{Text: "이 문서는 `아메리카노` 이다", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if a.Subject != "아메리카노" {
					t.Fatalf("subject = %q", a.Subject)
				}
			},
		},
		{
			name:  "See",
			code:  codebook.SectionSee,
			items: []scan.Item{{Text: "`a` 는 `b` 에서 참조한다", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Sees) != 1 {
					t.Fatalf("sees = %+v", a.Sees)
				}
			},
		},
		{
			name:  "Definitions",
			code:  codebook.SectionDefinitions,
			items: []scan.Item{{Text: "`문턱` 은 10 을 뜻한다", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Defs) != 1 {
					t.Fatalf("defs = %+v", a.Defs)
				}
			},
		},
		{
			name:  "Rules",
			code:  codebook.SectionRules,
			items: []scan.Item{{Text: "`신용점수` 가 650 이상이면 `신용 충분` 이다", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Rules) != 1 {
					t.Fatalf("rules = %+v", a.Rules)
				}
			},
		},
		{
			name:  "Cases",
			code:  codebook.SectionCases,
			items: []scan.Item{{Text: "`컵 놓기 가능` 의 경우에", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Cases) != 1 {
					t.Fatalf("cases = %+v", a.Cases)
				}
			},
		},
		{
			name:  "Provides",
			code:  codebook.SectionProvides,
			items: []scan.Item{{Text: "`X` 를 한다면", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Provides) != 1 {
					t.Fatalf("provides = %+v", a.Provides)
				}
			},
		},
		{
			name:  "DefaultInternal",
			code:  codebook.SectionInternal,
			items: []scan.Item{{Text: "`버튼 눌림` 발생 시", No: 1}},
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Internals) != 1 {
					t.Fatalf("internals = %+v", a.Internals)
				}
			},
		},
		{
			name:  "DefaultUnmapped",
			code:  codebook.SectionCode(99),
			items: nil,
			check: func(t *testing.T, a *ast.Document) {
				if len(a.Internals) != 0 {
					t.Fatalf("internals = %+v", a.Internals)
				}
			},
		},
	}

	for _, c := range cases {
		a := &ast.Document{}
		anchors, err := sectionAnchorsKo(a, c.code, c.items, "t.md")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", c.name, err)
		}
		if c.items != nil && len(anchors) == 0 {
			t.Fatalf("%s: expected non-empty anchors", c.name)
		}
		c.check(t, a)
	}
}
