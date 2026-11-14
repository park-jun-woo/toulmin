//ff:func feature=tangl type=parser control=sequence
//ff:what TestIndexCaseRefs — case slots index by line; a supersede keeps an occupied line unless it heads a chain
package en

import (
	"reflect"
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestIndexCaseRefs(t *testing.T) {
	c := ast.Case{
		Line:     10,
		Requires: []ast.Require{{Line: 11}},
		Nodes:    []ast.Node{{Name: "A", Line: 12}, {Name: "B", Line: 13}},
		Attacks:  []ast.Attack{{Line: 14}},
		Execs:    []ast.Exec{{Line: 15}},
		Supersedes: []ast.Supersede{
			{Superseder: "A", Superseded: "B", Line: 16},
			{Superseder: "A", Superseded: "B", Line: 12},
			{Superseder: "A", Superseded: "B", Line: 13, Chain: []string{"A", "B", "C"}},
		},
	}
	idx := map[int]doc.StmtRef{99: {Sec: toulmin.SectionSee}}
	indexCaseRefs(idx, 2, c)
	ref := func(path ...int) doc.StmtRef { return doc.StmtRef{Sec: toulmin.SectionCases, Path: path} }
	want := map[int]doc.StmtRef{
		99: {Sec: toulmin.SectionSee},
		10: ref(2),
		11: ref(2, 0, 0),
		12: ref(2, 1, 0),
		13: ref(2, 4, 2),
		14: ref(2, 2, 0),
		15: ref(2, 3, 0),
		16: ref(2, 4, 0),
	}
	if !reflect.DeepEqual(idx, want) {
		t.Fatalf("idx = %+v\nwant %+v", idx, want)
	}
}
