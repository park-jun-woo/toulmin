//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestResolveInternalSlot — both slotKinds resolve valid indices and reject negative/out-of-range/unknown (0 included)
package doc

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestResolveInternalSlot verifies resolveInternalSlot's two Internal slot
// kinds (Run=1, Check=2): each resolves a valid in-range index and rejects
// both a negative index and an index at/past len; slotKind 0 (unused for
// Internal) and any other unrecognized value hit the default branch.
func TestResolveInternalSlot(t *testing.T) {
	in := &ast.Internal{
		Runs:   []string{"ir"},
		Checks: []string{"ic"},
	}

	cases := []struct {
		name     string
		slotKind int
		slotIdx  int
		want     any
		wantErr  bool
	}{
		{"run valid", 1, 0, in.Runs[0], false},
		{"run negative", 1, -1, nil, true},
		{"run oob", 1, 1, nil, true},
		{"check valid", 2, 0, in.Checks[0], false},
		{"check negative", 2, -1, nil, true},
		{"check oob", 2, 1, nil, true},
		{"slotKind 0 unused", 0, 0, nil, true},
		{"unknown slotKind", 3, 0, nil, true},
	}

	for _, tc := range cases {
		got, err := resolveInternalSlot(in, tc.slotKind, tc.slotIdx)
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
