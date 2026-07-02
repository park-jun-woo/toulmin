//ff:func feature=tangl type=codegen control=sequence
//ff:what TestWriteImportBlock — tests writeImportBlock for aliased and unaliased import branches
package gen

import (
	"strings"
	"testing"
)

func TestWriteImportBlock(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var w strings.Builder
		writeImportBlock(&w, nil)
		got := w.String()
		want := "import (\n)\n"
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("aliased and unaliased", func(t *testing.T) {
		var w strings.Builder
		specs := []importSpec{
			{Path: "fmt"},
			{Alias: "toulmin", Path: "github.com/park-jun-woo/toulmin/pkg/toulmin"},
		}
		writeImportBlock(&w, specs)
		got := w.String()
		if !strings.Contains(got, "\t\"fmt\"\n") {
			t.Errorf("expected unaliased fmt import, got %q", got)
		}
		if !strings.Contains(got, "\ttoulmin \"github.com/park-jun-woo/toulmin/pkg/toulmin\"\n") {
			t.Errorf("expected aliased toulmin import, got %q", got)
		}
		if !strings.HasPrefix(got, "import (\n") || !strings.HasSuffix(got, ")\n") {
			t.Errorf("expected import block wrapper, got %q", got)
		}
	})
}
