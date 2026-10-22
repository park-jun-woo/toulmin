//ff:func feature=tangl type=parser control=sequence
//ff:what TestBuildRefIndex — index a synthetic AST across all line-bearing sections
package en

import (
	"reflect"
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestBuildRefIndex builds an AST with one element of each line-bearing kind at a
// distinct line and asserts the resulting line→StmtRef map matches the README §2
// path table (Cases slotKinds, Definition fields, endpoint requires). Subject and
// the run/check bodies carry no Line and must be absent.
func TestBuildRefIndex(t *testing.T) {
	a := &ast.Document{
		Subject: "s",
		Sees:    []ast.See{{Alias: "p", Line: 10}},
		Defs: []ast.Definition{{
			Name:   "d",
			Line:   20,
			Fields: []ast.Field{{Name: "f", Line: 21}},
		}},
		Rules: []ast.InlineRule{{Name: "r", Line: 30}},
		Cases: []ast.Case{{
			Name:     "c",
			Line:     40,
			Requires: []ast.Require{{Field: "q", Line: 41}},
			Nodes:    []ast.Node{{Name: "n", Line: 42}},
			Attacks:  []ast.Attack{{Target: "n", Attacker: "n", Line: 43}},
			Execs:    []ast.Exec{{Node: "n", Line: 44}},
			Supersedes: []ast.Supersede{
				{Superseder: "n", Superseded: "n", Line: 45},
			},
		}},
		Provides: []ast.Endpoint{{
			Name:     "ep",
			Line:     50,
			Requires: []ast.Require{{Field: "q", Line: 51}},
			Runs:     []string{"c"},
		}},
		Internals: []ast.Internal{{Line: 60, Runs: []string{"c"}}},
	}

	got := buildRefIndex(a)
	want := map[int]doc.StmtRef{
		10: {Sec: toulmin.SectionSee, Path: []int{0}},
		20: {Sec: toulmin.SectionDefinitions, Path: []int{0}},
		21: {Sec: toulmin.SectionDefinitions, Path: []int{0, 0}},
		30: {Sec: toulmin.SectionRules, Path: []int{0}},
		40: {Sec: toulmin.SectionCases, Path: []int{0}},
		41: {Sec: toulmin.SectionCases, Path: []int{0, 0, 0}},
		42: {Sec: toulmin.SectionCases, Path: []int{0, 1, 0}},
		43: {Sec: toulmin.SectionCases, Path: []int{0, 2, 0}},
		44: {Sec: toulmin.SectionCases, Path: []int{0, 3, 0}},
		45: {Sec: toulmin.SectionCases, Path: []int{0, 4, 0}},
		50: {Sec: toulmin.SectionProvides, Path: []int{0}},
		51: {Sec: toulmin.SectionProvides, Path: []int{0, 0, 0}},
		60: {Sec: toulmin.SectionInternal, Path: []int{0}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildRefIndex mismatch\n got %v\nwant %v", got, want)
	}
}
