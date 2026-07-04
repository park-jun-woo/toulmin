//ff:func feature=tangl type=codegen control=sequence
//ff:what TestFindModulePath — covers success, missing go.mod, unreadable go.mod, and no module line
package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindModulePath(t *testing.T) {
	t.Run("reads the declared module path", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/m\n\ngo 1.18\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		gotRoot, module, err := findModulePath(root)
		if err != nil {
			t.Fatal(err)
		}
		if module != "example.com/m" {
			t.Errorf("module = %q, want example.com/m", module)
		}
		if _, statErr := os.Stat(filepath.Join(gotRoot, "go.mod")); statErr != nil {
			t.Errorf("returned root %q has no go.mod", gotRoot)
		}
	})

	t.Run("no enclosing go.mod errors", func(t *testing.T) {
		if _, _, err := findModulePath(t.TempDir()); err == nil {
			t.Fatal("expected an error when no go.mod is found")
		}
	})

	t.Run("unreadable go.mod errors", func(t *testing.T) {
		// A go.mod that is a directory passes the Stat existence check in
		// findGoModDir but fails os.ReadFile.
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, "go.mod"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, _, err := findModulePath(root); err == nil {
			t.Fatal("expected a read error for a directory-shaped go.mod")
		}
	})

	t.Run("go.mod without a module line errors", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("go 1.18\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, _, err := findModulePath(root)
		if err == nil || !strings.Contains(err.Error(), "no module line") {
			t.Fatalf("expected a no-module-line error, got %v", err)
		}
	})
}
