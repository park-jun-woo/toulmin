//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGenerateGrounded — generates a grounded-clause fixture and checks gofmt-clean, deterministic output with the Grounds API call
package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/parser"
)

// TestGenerateGrounded parses a fixture that exercises the grounded clause
// (single and multi-name, the latter combined with a level), then checks the
// generated code re-formats cleanly, contains the Grounds API call in document
// order, and is deterministic across repeated parse+generate runs.
func TestGenerateGrounded(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/grounded.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	// Determinism: a second parse must yield byte-identical Grounds slices.
	doc2, err := parser.Parse("../parser/testdata/grounded.md")
	if err != nil {
		t.Fatalf("Parse (second): %v", err)
	}
	if fp1, fp2 := groundsFingerprint(doc), groundsFingerprint(doc2); fp1 != fp2 {
		t.Fatalf("non-deterministic Grounds: %q vs %q", fp1, fp2)
	}

	src, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("re-format.Source: %v\n%s", err, src)
	}
	src2, err := Generate(doc2)
	if err != nil {
		t.Fatalf("Generate (second): %v", err)
	}
	if src != src2 {
		t.Fatalf("non-deterministic generated output:\n--- first ---\n%s\n--- second ---\n%s", src, src2)
	}
	for _, want := range []string{
		`.Grounds("source-body")`,
		`.Grounds("source-body", "mx")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated code missing %q; got:\n%s", want, src)
		}
	}
}
