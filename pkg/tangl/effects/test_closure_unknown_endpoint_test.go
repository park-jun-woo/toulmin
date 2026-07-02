//ff:func feature=tangl type=analyzer control=sequence
//ff:what TestClosureUnknownEndpoint — an endpoint name absent from Provides is an error
package effects

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestClosureUnknownEndpoint verifies that Closure rejects an endpoint name
// that is not declared in the document's tangl:Provides section.
func TestClosureUnknownEndpoint(t *testing.T) {
	doc, err := parser.Parse("../parser/testdata/transfer.md")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if _, err := Closure(doc, "nope"); err == nil {
		t.Fatal("expected error for unknown endpoint, got nil")
	}
}
