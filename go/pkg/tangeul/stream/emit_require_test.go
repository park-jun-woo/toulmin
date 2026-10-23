//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitRequire — encodes a required field with and without an `as Type` label
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitRequire covers the r.Type != "" branch (an extra `as Type` label) in
// both directions. appendNode/appendLabel/appendEdge never fail here (NodeField
// and PropRequires are always emitted with consistent flags), so their
// error-return lines are unreachable through any legitimate call and are not
// exercised.
func TestEmitRequire(t *testing.T) {
	span := toulmin.Span{}
	cases := []struct {
		name string
		r    ast.Require
	}{
		{"no type", ast.Require{Field: "amount"}},
		{"with type", ast.Require{Field: "amount", Type: "Currency"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e := &enc{tt: &tidTable{byRef: map[string]toulmin.TID{}}, wire: testWire, next: 100}
			got, err := emitRequire(e, c.r, 5, 4, span)
			if err != nil {
				t.Fatalf("emitRequire(%+v): %v", c.r, err)
			}
			if len(got) == 0 {
				t.Errorf("emitRequire(%+v): expected non-empty words", c.r)
			}
		})
	}
}
