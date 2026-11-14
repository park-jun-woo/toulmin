//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitSupersede — undeclared superseder, undeclared superseded, and a successful exclusion edge
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitSupersede drives emitSupersede's two undeclared-node hard errors
// (superseder, then superseded) and the success path where both resolve through
// nodeByCase.
func TestEmitSupersede(t *testing.T) {
	e := testEnc(&tidTable{
		nodeByCase: map[string]toulmin.TID{
			nodeKey(0, "a"): 7,
			nodeKey(0, "b"): 8,
		},
	})
	ref := doc.StmtRef{Sec: 5, Path: []int{0, 4, 0}}
	span := toulmin.Span{Off: 1, Len: 2}

	cases := []struct {
		name    string
		s       ast.Supersede
		wantErr bool
	}{
		{"undeclared superseder", ast.Supersede{Superseder: "missing", Superseded: "b"}, true},
		{"undeclared superseded", ast.Supersede{Superseder: "a", Superseded: "missing"}, true},
		{"success", ast.Supersede{Superseder: "a", Superseded: "b"}, false},
	}
	for _, c := range cases {
		got, err := emitSupersede(e, ref, c.s, span)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: emitSupersede: %v", c.name, err)
		}
		if len(got) == 0 {
			t.Errorf("%s: expected non-empty words", c.name)
		}
	}
}
