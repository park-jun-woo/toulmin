//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteRequiredGuard — tests writeRequiredGuard for empty-fields, returnsResults, and single-error-return branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteRequiredGuard(t *testing.T) {
	t.Run("empty fields", func(t *testing.T) {
		var w strings.Builder
		writeRequiredGuard(&w, nil, true)
		if w.String() != "" {
			t.Errorf("expected empty output, got %q", w.String())
		}
	})

	t.Run("returns results", func(t *testing.T) {
		var w strings.Builder
		writeRequiredGuard(&w, []string{"amount", "status"}, true)
		got := w.String()
		if !strings.Contains(got, `tangl.Required(ctx, "amount", "status")`) {
			t.Errorf("expected required call, got %q", got)
		}
		if !strings.Contains(got, "return nil, err") {
			t.Errorf("expected nil,err return, got %q", got)
		}
	})

	t.Run("not returns results", func(t *testing.T) {
		var w strings.Builder
		writeRequiredGuard(&w, []string{"amount"}, false)
		got := w.String()
		if !strings.Contains(got, `tangl.Required(ctx, "amount")`) {
			t.Errorf("expected required call, got %q", got)
		}
		if !strings.Contains(got, "\t\treturn err\n") {
			t.Errorf("expected err return, got %q", got)
		}
		if strings.Contains(got, "return nil, err") {
			t.Errorf("unexpected nil,err return, got %q", got)
		}
	})
}
