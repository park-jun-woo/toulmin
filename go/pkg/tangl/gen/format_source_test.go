//ff:func feature=tangl type=codegen control=sequence
//ff:what TestFormatSource — tests formatSource for valid-source success and invalid-source error branches
package gen

import (
	"strings"
	"testing"
)

func TestFormatSource(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		src := "package p\nfunc f(){\nreturn\n}\n"
		out, err := formatSource(src)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(out, "package p") {
			t.Errorf("expected formatted output to contain package clause, got %q", out)
		}
	})

	t.Run("error", func(t *testing.T) {
		src := "package p\nfunc f( {\n"
		out, err := formatSource(src)
		if err == nil {
			t.Fatal("expected error for invalid Go source")
		}
		if out != "" {
			t.Errorf("expected empty output on error, got %q", out)
		}
	})
}
