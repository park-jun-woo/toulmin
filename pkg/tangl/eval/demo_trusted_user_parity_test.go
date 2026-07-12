//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what TestDemoTrustedUserParity — the field-path/spec expression matches moderate.IsTrustedUser verdicts
package eval

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/moderate"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// TestDemoTrustedUserParity rewrites the B-group representative demo
// predicate moderate.IsTrustedUser (pkg/moderate/rule_is_trusted_user.go,
// `a.TrustScore >= tb.MinScore`) as `author.trustScore is at least
// spec TrustScoreSpec.MinScore` and checks the eval-compiled closure returns
// the same bool as the hand-written Go predicate, including the missing and
// wrong-typed author cases.
func TestDemoTrustedUserParity(t *testing.T) {
	cond := ast.Compare{
		Lhs: ast.Path{Steps: []string{"author", "trustScore"}},
		Op:  "is at least",
		Rhs: ast.SpecRef{Steps: []string{"TrustScoreSpec", "MinScore"}},
	}
	fn, err := Compile(ast.InlineRule{Name: "trusted author", Cond: cond}, nil)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	specs := toulmin.Specs{&moderate.TrustScoreSpec{MinScore: 0.7}}
	cases := []struct {
		name   string
		author any
	}{
		{"trusted", &moderate.Author{ID: "a", TrustScore: 0.9}},
		{"boundary", &moderate.Author{ID: "b", TrustScore: 0.7}},
		{"untrusted", &moderate.Author{ID: "c", TrustScore: 0.5}},
		{"wrong type", "nobody"},
	}
	for _, tc := range cases {
		ctx := toulmin.NewContext()
		ctx.Set("author", tc.author)
		got, _ := fn(ctx, specs)
		want, _ := moderate.IsTrustedUser(ctx, specs)
		if got != want {
			t.Errorf("%s: eval=%t, hand-written=%t", tc.name, got, want)
		}
	}
}
