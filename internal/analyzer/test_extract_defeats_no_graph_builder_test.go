//ff:func feature=analyzer type=analyzer control=sequence
//ff:what TestExtractDefeatsNoGraphBuilder — tests extraction when no Graph builder exists
package analyzer

import (
	"testing"
)

func TestExtractDefeatsNoGraphBuilder(t *testing.T) {
	path := writeGoFile(t, `package example

func Foo() {}
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 0 {
		t.Errorf("expected 0 graphs, got %d", len(graphs))
	}
}
