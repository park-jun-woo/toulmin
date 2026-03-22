//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphValid — tests validation passes for a valid defeats graph
package graph

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestValidateGraphValid(t *testing.T) {
	metas := []toulmin.RuleMeta{
		{Name: "W", Strength: toulmin.Defeasible},
		{Name: "R", Strength: toulmin.Defeasible, Defeats: []string{"W"}},
	}
	if err := ValidateGraph(metas); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
