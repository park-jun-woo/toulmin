//ff:func feature=cli type=command control=sequence
//ff:what TestFindDocumentBySubjectMatch — findDocumentBySubject skips subject mismatches, returns a single match, and errors on an ambiguous match
package tanglcli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFindDocumentBySubjectMatch covers the matching branches of
// findDocumentBySubject: a subject mismatch is skipped, a single match is
// returned, and an ambiguous match across two documents is an error.
func TestFindDocumentBySubjectMatch(t *testing.T) {
	t.Run("subject mismatch is skipped", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "beta.md"), []byte(validTanglDoc("beta")), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		doc, err := findDocumentBySubject(dir, "alpha")
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

	t.Run("single match returns the document", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "gamma.md"), []byte(validTanglDoc("gamma")), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		doc, err := findDocumentBySubject(dir, "gamma")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc == nil {
			t.Fatal("doc = nil, want a match")
		}
		if doc.Subject != "gamma" {
			t.Errorf("doc.Subject = %q, want %q", doc.Subject, "gamma")
		}
	})

	t.Run("ambiguous subject across two documents is an error", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "doc_a.md"), []byte(validTanglDoc("delta")), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "doc_b.md"), []byte(validTanglDoc("delta")), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		doc, err := findDocumentBySubject(dir, "delta")
		if err == nil {
			t.Fatal("expected ambiguous-subject error, got nil")
		}
		if !strings.Contains(err.Error(), "ambiguous subject") {
			t.Errorf("err = %v, want 'ambiguous subject' message", err)
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
	})
}
