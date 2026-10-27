//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAppendRequire — routes a require to its Case or Endpoint owner, or errors on an unknown owner type
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestAppendRequire covers all three arms of appendRequire's switch on the
// owner's recorded NodeType: NodeCase appends to the case's Requires,
// NodeEndpoint appends to the endpoint's Requires, and any other (including
// unset/zero) NodeType is a hard error.
func TestAppendRequire(t *testing.T) {
	newDec := func() *dec {
		return &dec{
			ast: &ast.Document{
				Cases:    []ast.Case{{Name: "c0"}},
				Provides: []ast.Endpoint{{Name: "ep0"}},
			},
			nodeType:  map[toulmin.TID]toulmin.NodeType{},
			caseByTID: map[toulmin.TID]int{},
			epByTID:   map[toulmin.TID]int{},
		}
	}

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "case owner",
			run: func(t *testing.T) {
				d := newDec()
				d.nodeType[4] = toulmin.NodeCase
				d.caseByTID[4] = 0
				ref, err := appendRequire(d, 4, ast.Require{Field: "amount"})
				if err != nil {
					t.Fatalf("appendRequire: %v", err)
				}
				if len(d.ast.Cases[0].Requires) != 1 || d.ast.Cases[0].Requires[0].Field != "amount" {
					t.Errorf("case requires = %v", d.ast.Cases[0].Requires)
				}
				if ref.Sec != toulmin.SectionCases || len(ref.Path) != 3 || ref.Path[0] != 0 || ref.Path[1] != 0 || ref.Path[2] != 0 {
					t.Errorf("ref = %+v", ref)
				}
			},
		},
		{
			name: "endpoint owner",
			run: func(t *testing.T) {
				d := newDec()
				d.nodeType[12] = toulmin.NodeEndpoint
				d.epByTID[12] = 0
				ref, err := appendRequire(d, 12, ast.Require{Field: "user"})
				if err != nil {
					t.Fatalf("appendRequire: %v", err)
				}
				if len(d.ast.Provides[0].Requires) != 1 || d.ast.Provides[0].Requires[0].Field != "user" {
					t.Errorf("endpoint requires = %v", d.ast.Provides[0].Requires)
				}
				if ref.Sec != toulmin.SectionProvides || len(ref.Path) != 3 || ref.Path[0] != 0 || ref.Path[1] != 0 || ref.Path[2] != 0 {
					t.Errorf("ref = %+v", ref)
				}
			},
		},
		{
			name: "unknown owner type errors",
			run: func(t *testing.T) {
				d := newDec()
				_, err := appendRequire(d, 99, ast.Require{Field: "x"})
				if err == nil {
					t.Fatalf("want error for unknown owner")
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
