//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_List — tests listing active feature flags
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_List(t *testing.T) {
	flags := NewFlags()

	g1 := toulmin.NewGraph("feature:a")
	g1.Rule(IsBetaUser)
	flags.Register("a", g1)

	g2 := toulmin.NewGraph("feature:b")
	g2.Rule(IsRegion).With(&RegionSpec{Region: "KR"})
	flags.Register("b", g2)

	ctx := &UserContext{
		Region:     "KR",
		Attributes: map[string]any{"beta": true},
	}
	active, err := flags.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(active) != 2 {
		t.Errorf("expected 2 active features, got %d", len(active))
	}
}
