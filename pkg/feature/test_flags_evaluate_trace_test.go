//ff:func feature=feature type=engine control=sequence
//ff:what TestFlags_EvaluateTrace — tests EvaluateTrace returns trace data
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestFlags_EvaluateTrace(t *testing.T) {
	flags := NewFlags()

	g := toulmin.NewGraph("feature:dark-mode")
	g.Rule(IsBetaUser)
	flags.Register("dark-mode", g)

	ctx := &UserContext{Attributes: map[string]any{"beta": true}}
	result, err := flags.EvaluateTrace("dark-mode", ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Enabled {
		t.Error("expected enabled")
	}
	if len(result.Trace) == 0 {
		t.Error("expected non-empty trace")
	}
}
