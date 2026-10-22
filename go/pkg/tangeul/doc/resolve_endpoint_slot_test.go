//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestResolveEndpointSlot — all three slotKinds resolve valid indices and reject negative/out-of-range/unknown
package doc

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestResolveEndpointSlot verifies resolveEndpointSlot's three Endpoint slot
// kinds (Require=0, Run=1, Check=2): each resolves a valid in-range index and
// rejects both a negative index and an index at/past len; an unrecognized
// slotKind hits the default branch.
func TestResolveEndpointSlot(t *testing.T) {
	ep := &ast.Endpoint{
		Name:     "ep0",
		Requires: []ast.Require{{Field: "rq"}},
		Runs:     []string{"ru"},
		Checks:   []string{"ch"},
	}

	cases := []struct {
		name     string
		slotKind int
		slotIdx  int
		want     any
		wantErr  bool
	}{
		{"require valid", 0, 0, ep.Requires[0], false},
		{"require negative", 0, -1, nil, true},
		{"require oob", 0, 1, nil, true},
		{"run valid", 1, 0, ep.Runs[0], false},
		{"run negative", 1, -1, nil, true},
		{"run oob", 1, 1, nil, true},
		{"check valid", 2, 0, ep.Checks[0], false},
		{"check negative", 2, -1, nil, true},
		{"check oob", 2, 1, nil, true},
		{"unknown slotKind", 3, 0, nil, true},
	}

	for _, tc := range cases {
		got, err := resolveEndpointSlot(ep, tc.slotKind, tc.slotIdx)
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
