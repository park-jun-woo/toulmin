//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_Percentage_Deterministic — tests deterministic percentage-based rollout
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_Percentage_Deterministic(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:rollout")
	g.Warrant(IsUserInPercentage, 0.5, 1.0)
	flags.Register("rollout", g)

	ctx := &UserContext{ID: "user-42"}
	r1, _ := flags.IsEnabled("rollout", ctx)
	r2, _ := flags.IsEnabled("rollout", ctx)
	if r1 != r2 {
		t.Error("expected deterministic result for same user")
	}
}
