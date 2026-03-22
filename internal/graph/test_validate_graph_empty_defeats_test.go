//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphEmptyDefeats — tests validation passes when no defeats exist
package graph

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestValidateGraphEmptyDefeats(t *testing.T) {
	metas := []toulmin.RuleMeta{
		{Name: "W", Strength: toulmin.Defeasible},
	}
	if err := ValidateGraph(metas); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
