//ff:func feature=cli type=command control=sequence
//ff:what TestFindDocumentBySubjectSkip — findDocumentBySubject skips walk errors, non-.tangeul files, unreadable files, and corrupt streams
package tanglcli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFindDocumentBySubjectSkip covers the walk-error and skip branches of
// findDocumentBySubject: a nonexistent root propagates the walk error;
// subdirectories, non-.tangeul files, and .md files are skipped; an unreadable
// .tangeul file is skipped; and a corrupt .tangeul (load error) is skipped.
func TestFindDocumentBySubjectSkip(t *testing.T) {
	t.Run("nonexistent root propagates walk error", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "does-not-exist")
		doc, err := findDocumentBySubject(dir, "whatever")
		if err == nil {
			t.Fatal("expected walk error, got nil")
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
	})

	t.Run("skips subdirectories, non-.tangeul files, and .md files", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
			t.Fatalf("Mkdir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "note.txt"), []byte("hello"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "policy.md"), []byte("## tangl:Subject\n- x\n"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		doc, err := findDocumentBySubject(dir, "whatever")
		if err == nil {
			t.Fatal("expected not-found error, got nil")
		}
		if !strings.Contains(err.Error(), "no document with subject") {
			t.Errorf("err = %v, want 'no document with subject' message", err)
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
	})

	t.Run("unreadable .tangeul file is skipped", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("running as root: file permissions are not enforced")
		}
		dir := t.TempDir()
		path := filepath.Join(dir, "secret.tangeul")
		saveTangeul(t, path, buildSubjectDoc("secret"))
		if err := os.Chmod(path, 0o000); err != nil {
			t.Fatalf("Chmod: %v", err)
		}
		defer os.Chmod(path, 0o644)
		doc, err := findDocumentBySubject(dir, "secret")
		if err == nil {
			t.Fatal("expected not-found error, got nil")
		}
		if !strings.Contains(err.Error(), "no document with subject") {
			t.Errorf("err = %v, want 'no document with subject' message", err)
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
	})

	t.Run("corrupt .tangeul (load error) is skipped", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "bad.tangeul"), []byte("not a stream"), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		doc, err := findDocumentBySubject(dir, "whatever")
		if err == nil {
			t.Fatal("expected not-found error, got nil")
		}
		if !strings.Contains(err.Error(), "no document with subject") {
			t.Errorf("err = %v, want 'no document with subject' message", err)
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
	})
}
