//ff:func feature=tangl type=codegen control=sequence
//ff:what TestFindGoModDir — covers upward discovery of go.mod and the no-module error
package gen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGoModDir(t *testing.T) {
	t.Run("finds enclosing module by walking up", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		sub := filepath.Join(root, "a", "b")
		if err := os.MkdirAll(sub, 0o755); err != nil {
			t.Fatal(err)
		}
		got, err := findGoModDir(sub)
		if err != nil {
			t.Fatal(err)
		}
		if _, statErr := os.Stat(filepath.Join(got, "go.mod")); statErr != nil {
			t.Fatalf("returned dir %q has no go.mod: %v", got, statErr)
		}
	})

	t.Run("no enclosing module errors at filesystem root", func(t *testing.T) {
		// A fresh temp dir under the system temp root has no go.mod anywhere
		// above it, so the walk reaches the filesystem root and errors.
		if _, err := findGoModDir(t.TempDir()); err == nil {
			t.Fatal("expected an error when no go.mod exists above the directory")
		}
	})
}
