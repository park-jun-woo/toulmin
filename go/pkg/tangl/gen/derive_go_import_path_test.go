//ff:func feature=tangl type=codegen control=sequence
//ff:what TestDeriveGoImportPath — tests module-root and subdirectory import path derivation via go.mod lookup
package gen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeriveGoImportPath(t *testing.T) {
	root := t.TempDir()
	gomod := "module example.com/reins\n\ngo 1.18\n"
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(root, "gate")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	t.Run("module root document", func(t *testing.T) {
		got, err := deriveGoImportPath(filepath.Join(root, "turn.geul"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "example.com/reins"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("subdirectory document", func(t *testing.T) {
		got, err := deriveGoImportPath(filepath.Join(sub, "gate.geul"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "example.com/reins/gate"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("no enclosing module errors", func(t *testing.T) {
		_, err := deriveGoImportPath("/nonexistent-module-root-xyz/deeper/doc.geul")
		if err == nil {
			t.Fatal("expected an error when no go.mod is found")
		}
	})

	t.Run("abs failure on removed working directory", func(t *testing.T) {
		// A relative document path forces filepath.Abs to consult the working
		// directory, which no longer exists.
		withRemovedCwd(t)
		if _, err := deriveGoImportPath("rel/doc.geul"); err == nil {
			t.Fatal("expected filepath.Abs to fail on a removed working directory")
		}
	})
}
