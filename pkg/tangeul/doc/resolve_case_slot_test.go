//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestResolveCaseSlot — all four slotKinds resolve valid indices and reject negative/out-of-range/unknown
package doc

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestResolveCaseSlot verifies resolveCaseSlot's four Case slot kinds
// (Require=0, Node=1, Attack=2, Exec=3): each resolves a valid in-range index
// and rejects both a negative index and an index at/past len; an unrecognized
// slotKind hits the default branch.
func TestResolveCaseSlot(t *testing.T) {
	c := &ast.Case{
		Name:     "c0",
		Requires: []ast.Require{{Field: "rq"}},
		Nodes:    []ast.Node{{Name: "n0"}},
		Attacks:  []ast.Attack{{Target: "t", Attacker: "u"}},
		Execs:    []ast.Exec{{Node: "n0"}},
	}

	cases := []struct {
		name     string
		slotKind int
		slotIdx  int
		want     any
		wantErr  bool
	}{
		{"require valid", 0, 0, c.Requires[0], false},
		{"require negative", 0, -1, nil, true},
		{"require oob", 0, 1, nil, true},
		{"node valid", 1, 0, c.Nodes[0], false},
		{"node negative", 1, -1, nil, true},
		{"node oob", 1, 1, nil, true},
		{"attack valid", 2, 0, c.Attacks[0], false},
		{"attack negative", 2, -1, nil, true},
		{"attack oob", 2, 1, nil, true},
		{"exec valid", 3, 0, c.Execs[0], false},
		{"exec negative", 3, -1, nil, true},
		{"exec oob", 3, 1, nil, true},
		{"unknown slotKind", 4, 0, nil, true},
	}

	for _, tc := range cases {
		got, err := resolveCaseSlot(c, tc.slotKind, tc.slotIdx)
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
}
