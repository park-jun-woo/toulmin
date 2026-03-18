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

func TestValidateGraphUnknownTarget(t *testing.T) {
	metas := []toulmin.RuleMeta{
		{Name: "W", Strength: toulmin.Defeasible},
		{Name: "R", Strength: toulmin.Defeasible, Defeats: []string{"Ghost"}},
	}
	if err := ValidateGraph(metas); err == nil {
		t.Fatal("expected error for unknown defeat target")
	}
}

func TestValidateGraphEmptyDefeats(t *testing.T) {
	metas := []toulmin.RuleMeta{
		{Name: "W", Strength: toulmin.Defeasible},
	}
	if err := ValidateGraph(metas); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
