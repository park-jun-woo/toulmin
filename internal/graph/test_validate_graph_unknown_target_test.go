//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphUnknownTarget — tests validation fails for unknown defeat target
package graph

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestValidateGraphUnknownTarget(t *testing.T) {
	metas := []toulmin.RuleMeta{
		{Name: "W", Strength: toulmin.Defeasible},
		{Name: "R", Strength: toulmin.Defeasible, Defeats: []string{"Ghost"}},
	}
	if err := ValidateGraph(metas); err == nil {
		t.Fatal("expected error for unknown defeat target")
	}
}
