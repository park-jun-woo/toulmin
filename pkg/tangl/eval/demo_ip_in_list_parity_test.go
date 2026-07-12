//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestDemoIPInListParity — the spec-sourced exists expression matches policy.IsIPInList verdicts
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/policy"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestDemoIPInListParity rewrites the C-group demo predicate
// policy.IsIPInList (pkg/policy/rule_is_ip_in_list.go) as
// `any of spec IPListSpec.List as ip satisfies ip equals clientIP` — an
// exists whose source is a spec reference and whose flat Field left side
// names the binder — and checks the eval-compiled closure returns the same
// bool as the hand-written Go predicate.
func TestDemoIPInListParity(t *testing.T) {
	cond := ast.Exists{
		Source: ast.SpecRef{Steps: []string{"IPListSpec", "List"}},
		Binder: "ip",
		Cond: ast.Compare{
			Field: "ip",
			Op:    "equals",
			Rhs:   ast.Path{Steps: []string{"clientIP"}},
		},
	}
	fn, err := Compile(ast.InlineRule{Name: "ip in list", Cond: cond}, nil)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	specs := toulmin.Specs{&policy.IPListSpec{Purpose: "blocklist", List: []string{"1.1.1.1", "2.2.2.2"}}}
	cases := []struct {
		name     string
		clientIP any
	}{
		{"listed", "2.2.2.2"},
		{"not listed", "9.9.9.9"},
		{"wrong type", 7},
	}
	for _, tc := range cases {
		ctx := toulmin.NewContext()
		ctx.Set("clientIP", tc.clientIP)
		got, _ := fn(ctx, specs)
		want, _ := policy.IsIPInList(ctx, specs)
		if got != want {
			t.Errorf("%s: eval=%t, hand-written=%t", tc.name, got, want)
		}
	}
}
