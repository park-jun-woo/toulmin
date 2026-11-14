//ff:func feature=tangl type=parser control=sequence
//ff:what TestParse — tests Parse for file-read-error and successful-parse branches
package parser

import "testing"

func TestParse(t *testing.T) {
	t.Run("FileReadError", func(t *testing.T) {
		_, err := Parse("testdata/does-not-exist.md")
		if err == nil {
			t.Fatal("expected an error for a nonexistent file")
		}
	})

	t.Run("Success", func(t *testing.T) {
		doc, err := Parse("testdata/access.md")
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if doc.Subject != "api" {
			t.Errorf("Subject = %q, want api", doc.Subject)
		}
	})
}
