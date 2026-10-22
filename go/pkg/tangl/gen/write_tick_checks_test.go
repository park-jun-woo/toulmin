//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteTickChecks — tests writeTickChecks for empty and populated checks branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteTickChecks(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeTickChecks(&w, nil)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("populated", func(t *testing.T) {
		var w strings.Builder
		writeTickChecks(&w, []string{"order flow", "payment done"})
		got := w.String()
		if !strings.Contains(got, "\t\t_, _ = orderFlowGraph.Evaluate(ctx)\n") {
			t.Errorf("expected orderFlow evaluate call, got %q", got)
		}
		if !strings.Contains(got, "\t\t_, _ = paymentDoneGraph.Evaluate(ctx)\n") {
			t.Errorf("expected paymentDone evaluate call, got %q", got)
		}
	})
}
