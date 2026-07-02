//ff:func feature=cli type=command control=sequence
//ff:what TestResolveEffectsTarget — resolveEffectsTarget dispatches between the file+endpoint and subject.endpoint argument forms
package tanglcli

import (
	"strings"
	"testing"
)

const transferFixture = "../../pkg/tangl/parser/testdata/transfer.md"
const transferFixtureDir = "../../pkg/tangl/parser/testdata"

// TestResolveEffectsTarget covers every branch of resolveEffectsTarget: the
// two-arg (file, endpoint) form both succeeding and failing to parse, and
// the one-arg (subject.endpoint) form failing to split, failing to find a
// matching document, and succeeding.
func TestResolveEffectsTarget(t *testing.T) {
	t.Run("two args: file parses and endpoint is returned", func(t *testing.T) {
		doc, endpoint, err := resolveEffectsTarget([]string{transferFixture, "`transfer`"}, ".")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc == nil {
			t.Fatal("doc = nil, want a parsed document")
		}
		if doc.Subject != "transfer" {
			t.Errorf("doc.Subject = %q, want %q", doc.Subject, "transfer")
		}
		if endpoint != "transfer" {
			t.Errorf("endpoint = %q, want %q (backticks trimmed)", endpoint, "transfer")
		}
	})

	t.Run("two args: file fails to parse", func(t *testing.T) {
		doc, endpoint, err := resolveEffectsTarget([]string{"does-not-exist.md", "transfer"}, ".")
		if err == nil {
			t.Fatal("expected a parse error, got nil")
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
		if endpoint != "" {
			t.Errorf("endpoint = %q, want empty", endpoint)
		}
	})

	t.Run("one arg: missing dot fails to split", func(t *testing.T) {
		doc, endpoint, err := resolveEffectsTarget([]string{"noDotHere"}, transferFixtureDir)
		if err == nil {
			t.Fatal("expected a split error, got nil")
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
		if endpoint != "" {
			t.Errorf("endpoint = %q, want empty", endpoint)
		}
	})

	t.Run("one arg: subject not found under dir", func(t *testing.T) {
		doc, endpoint, err := resolveEffectsTarget([]string{"nonexistent-subject.endpoint"}, transferFixtureDir)
		if err == nil {
			t.Fatal("expected a not-found error, got nil")
		}
		if !strings.Contains(err.Error(), "no document with subject") {
			t.Errorf("err = %v, want 'no document with subject' message", err)
		}
		if doc != nil {
			t.Errorf("doc = %v, want nil", doc)
		}
		if endpoint != "" {
			t.Errorf("endpoint = %q, want empty", endpoint)
		}
	})

	t.Run("one arg: subject found under dir", func(t *testing.T) {
		doc, endpoint, err := resolveEffectsTarget([]string{"transfer.`transfer`"}, transferFixtureDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc == nil {
			t.Fatal("doc = nil, want a matched document")
		}
		if doc.Subject != "transfer" {
			t.Errorf("doc.Subject = %q, want %q", doc.Subject, "transfer")
		}
		if endpoint != "transfer" {
			t.Errorf("endpoint = %q, want %q (backticks trimmed)", endpoint, "transfer")
		}
	})
}
