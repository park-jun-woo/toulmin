//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteTickRunSequence — tests writeTickRunSequence for empty and populated runs branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteTickRunSequence(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeTickRunSequence(&w, nil)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("populated", func(t *testing.T) {
		var w strings.Builder
		writeTickRunSequence(&w, []string{"order flow", "payment done"})
		got := w.String()
		if !strings.Contains(got, "orderFlowGraph.Run(ctx)") {
			t.Errorf("expected orderFlow run call, got %q", got)
		}
		if !strings.Contains(got, "paymentDoneGraph.Run(ctx)") {
			t.Errorf("expected paymentDone run call, got %q", got)
		}
		if !strings.Contains(got, "tangl.Compensate(ctx)") {
			t.Errorf("expected compensate call, got %q", got)
		}
		if !strings.Contains(got, "\t\t\tcontinue\n") {
			t.Errorf("expected continue statement, got %q", got)
		}
	})
}
