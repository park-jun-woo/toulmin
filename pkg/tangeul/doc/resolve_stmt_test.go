//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestResolveStmt — ResolveStmt resolves all seven sections and rejects out-of-range paths
package doc

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestResolveStmt verifies ResolveStmt against the README §2 path table for all
// seven sections (Subject/See/Definitions/Rules/Cases/Provides/Internal, both the
// declaration and every body slot kind) and that malformed or out-of-range Paths
// are hard errors.
func TestResolveStmt(t *testing.T) {
	a := &ast.Document{
		Subject:   "sub",
		Sees:      []ast.See{{Alias: "a", Package: "p"}},
		Defs:      []ast.Definition{{Name: "d0", Fields: []ast.Field{{Name: "f0", Type: "T"}}}},
		Rules:     []ast.InlineRule{{Name: "r0"}},
		Cases:     []ast.Case{{Name: "c0", Requires: []ast.Require{{Field: "rq"}}, Nodes: []ast.Node{{Name: "n0"}}, Attacks: []ast.Attack{{Target: "t", Attacker: "u"}}, Execs: []ast.Exec{{Node: "n0"}}}},
		Provides:  []ast.Endpoint{{Name: "ep0", Requires: []ast.Require{{Field: "pr"}}, Runs: []string{"ru"}, Checks: []string{"ch"}}},
		Internals: []ast.Internal{{Runs: []string{"ir"}, Checks: []string{"ic"}}},
	}
	d := &Doc{AST: a}

	cases := []struct {
		name    string
		ref     StmtRef
		want    any
		wantErr bool
	}{
		{"subject", StmtRef{codebook.SectionSubject, nil}, "sub", false},
		{"see", StmtRef{codebook.SectionSee, []int{0}}, a.Sees[0], false},
		{"def", StmtRef{codebook.SectionDefinitions, []int{0}}, a.Defs[0], false},
		{"def field", StmtRef{codebook.SectionDefinitions, []int{0, 0}}, a.Defs[0].Fields[0], false},
		{"rule", StmtRef{codebook.SectionRules, []int{0}}, a.Rules[0], false},
		{"case", StmtRef{codebook.SectionCases, []int{0}}, a.Cases[0], false},
		{"case require", StmtRef{codebook.SectionCases, []int{0, 0, 0}}, a.Cases[0].Requires[0], false},
		{"case node", StmtRef{codebook.SectionCases, []int{0, 1, 0}}, a.Cases[0].Nodes[0], false},
		{"case attack", StmtRef{codebook.SectionCases, []int{0, 2, 0}}, a.Cases[0].Attacks[0], false},
		{"case exec", StmtRef{codebook.SectionCases, []int{0, 3, 0}}, a.Cases[0].Execs[0], false},
		{"endpoint", StmtRef{codebook.SectionProvides, []int{0}}, a.Provides[0], false},
		{"endpoint require", StmtRef{codebook.SectionProvides, []int{0, 0, 0}}, a.Provides[0].Requires[0], false},
		{"endpoint run", StmtRef{codebook.SectionProvides, []int{0, 1, 0}}, "ru", false},
		{"endpoint check", StmtRef{codebook.SectionProvides, []int{0, 2, 0}}, "ch", false},
		{"internal", StmtRef{codebook.SectionInternal, []int{0}}, a.Internals[0], false},
		{"internal run", StmtRef{codebook.SectionInternal, []int{0, 1, 0}}, "ir", false},
		{"internal check", StmtRef{codebook.SectionInternal, []int{0, 2, 0}}, "ic", false},
		{"subject bad path", StmtRef{codebook.SectionSubject, []int{0}}, nil, true},
		{"see bad len", StmtRef{codebook.SectionSee, []int{0, 0}}, nil, true},
		{"see oob", StmtRef{codebook.SectionSee, []int{5}}, nil, true},
		{"def bad len", StmtRef{codebook.SectionDefinitions, []int{0, 0, 0}}, nil, true},
		{"def index oob", StmtRef{codebook.SectionDefinitions, []int{9}}, nil, true},
		{"def field oob", StmtRef{codebook.SectionDefinitions, []int{0, 9}}, nil, true},
		{"rule bad len", StmtRef{codebook.SectionRules, []int{0, 0}}, nil, true},
		{"rule index oob", StmtRef{codebook.SectionRules, []int{9}}, nil, true},
		{"case bad len", StmtRef{codebook.SectionCases, []int{0, 9}}, nil, true},
		{"case index oob", StmtRef{codebook.SectionCases, []int{9}}, nil, true},
		{"case node oob", StmtRef{codebook.SectionCases, []int{0, 1, 9}}, nil, true},
		{"case bad slotkind", StmtRef{codebook.SectionCases, []int{0, 9, 0}}, nil, true},
		{"endpoint bad len", StmtRef{codebook.SectionProvides, []int{0, 0}}, nil, true},
		{"endpoint index oob", StmtRef{codebook.SectionProvides, []int{9}}, nil, true},
		{"internal bad len", StmtRef{codebook.SectionInternal, []int{0, 0}}, nil, true},
		{"internal index oob", StmtRef{codebook.SectionInternal, []int{9}}, nil, true},
		{"internal slotkind0", StmtRef{codebook.SectionInternal, []int{0, 0, 0}}, nil, true},
		{"unknown section", StmtRef{codebook.SectionCode(99), nil}, nil, true},
	}

	for _, tc := range cases {
		got, err := ResolveStmt(d, tc.ref)
		if tc.wantErr && err == nil {
			t.Errorf("%s: expected error, got %v", tc.name, got)
		}
		if tc.wantErr {
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %#v, want %#v", tc.name, got, tc.want)
		}
	}

	if _, err := ResolveStmt(nil, StmtRef{codebook.SectionSubject, nil}); err == nil {
		t.Errorf("nil Doc: expected error")
	}
}
