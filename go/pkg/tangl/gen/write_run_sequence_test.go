//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteRunSequence — tests writeRunSequence for empty/non-empty runs and returnsResults branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteRunSequence(t *testing.T) {
	t.Run("empty returns results", func(t *testing.T) {
		var w strings.Builder
		writeRunSequence(&w, nil, true)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("empty not returns results", func(t *testing.T) {
		var w strings.Builder
		writeRunSequence(&w, nil, false)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("returns results", func(t *testing.T) {
		var w strings.Builder
		writeRunSequence(&w, []string{"order flow"}, true)
		got := w.String()
		if !strings.Contains(got, "orderFlowGraph.Run(ctx)") {
			t.Errorf("expected run call, got %q", got)
		}
		if !strings.Contains(got, "return nil, tangl.Review(ctx, err, cerr)") {
			t.Errorf("expected review return with nil, got %q", got)
		}
		if !strings.Contains(got, "\t\treturn nil, err\n") {
			t.Errorf("expected zero return nil,err, got %q", got)
		}
	})

	t.Run("not returns results", func(t *testing.T) {
		var w strings.Builder
		writeRunSequence(&w, []string{"payment done", "shipping ready"}, false)
		got := w.String()
		if !strings.Contains(got, "paymentDoneGraph.Run(ctx)") {
			t.Errorf("expected first run call, got %q", got)
		}
		if !strings.Contains(got, "shippingReadyGraph.Run(ctx)") {
			t.Errorf("expected second run call, got %q", got)
		}
		if !strings.Contains(got, "\t\t\treturn tangl.Review(ctx, err, cerr)\n") {
			t.Errorf("expected review return without nil, got %q", got)
		}
		if !strings.Contains(got, "\t\treturn err\n") {
			t.Errorf("expected zero return err, got %q", got)
		}
		if strings.Contains(got, "return nil, err") {
			t.Errorf("unexpected nil,err return, got %q", got)
		}
	})
}
