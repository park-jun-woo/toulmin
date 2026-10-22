//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteExtraChecks — tests writeExtraChecks for no-results, empty-checks, and populated-checks branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteExtraChecks(t *testing.T) {
	t.Run("not returns results", func(t *testing.T) {
		var w strings.Builder
		writeExtraChecks(&w, []string{"some check"}, false)
		got := w.String()
		if got != "\treturn nil\n" {
			t.Errorf("expected simple return nil, got %q", got)
		}
	})

	t.Run("returns results no checks", func(t *testing.T) {
		var w strings.Builder
		writeExtraChecks(&w, nil, true)
		got := w.String()
		if !strings.Contains(got, "var out []toulmin.EvalResult") {
			t.Errorf("expected out var declaration, got %q", got)
		}
		if !strings.Contains(got, "return out, nil") {
			t.Errorf("expected final return, got %q", got)
		}
		if strings.Contains(got, "Graph.Evaluate") {
			t.Errorf("expected no loop body for empty checks, got %q", got)
		}
	})

	t.Run("returns results with checks", func(t *testing.T) {
		var w strings.Builder
		writeExtraChecks(&w, []string{"order received", "payment done"}, true)
		got := w.String()
		if !strings.Contains(got, "r0, err := orderReceivedGraph.Evaluate(ctx)") {
			t.Errorf("expected r0 evaluate call, got %q", got)
		}
		if !strings.Contains(got, "r1, err := paymentDoneGraph.Evaluate(ctx)") {
			t.Errorf("expected r1 evaluate call, got %q", got)
		}
		if !strings.Contains(got, "out = append(out, r0...)") {
			t.Errorf("expected r0 append, got %q", got)
		}
		if !strings.Contains(got, "out = append(out, r1...)") {
			t.Errorf("expected r1 append, got %q", got)
		}
		if !strings.Contains(got, "return out, nil") {
			t.Errorf("expected final return, got %q", got)
		}
	})
}
