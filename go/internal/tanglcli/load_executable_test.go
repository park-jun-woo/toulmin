//ff:func feature=cli type=command control=sequence
//ff:what TestLoadExecutable — a .md path is rejected with encode guidance, and any other extension is delegated to tangeul.Load (success and error)
package tanglcli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadExecutable covers loadExecutable's single branch: a .md path is
// rejected up front with a message pointing at `tangl encode`, while a non-.md
// path is handed to tangeul.Load, whose own outcome (success or error) is
// returned unchanged.
func TestLoadExecutable(t *testing.T) {
	t.Run("md extension is rejected", func(t *testing.T) {
		_, err := loadExecutable("subject.md")
		if err == nil || !strings.Contains(err.Error(), "text surface is not executable; run 'tangl encode subject.md' first") {
			t.Fatalf("err = %v, want text-surface guidance", err)
		}
	})

	t.Run("non-md path delegates to tangeul.Load on success", func(t *testing.T) {
		d := buildTransferDoc()
		path := filepath.Join(t.TempDir(), "transfer.geul")
		saveTangeul(t, path, d)

		got, err := loadExecutable(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.AST == nil || got.AST.Subject != "transfer" {
			t.Fatalf("got = %+v, want a loaded Doc with Subject %q", got, "transfer")
		}
	})

	t.Run("non-md path delegates to tangeul.Load on error", func(t *testing.T) {
		_, err := loadExecutable(filepath.Join(t.TempDir(), "does-not-exist.geul"))
		if err == nil {
			t.Fatal("expected a load error, got nil")
		}
	})
}
