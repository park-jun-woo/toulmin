//ff:func feature=analyzer type=analyzer control=sequence
//ff:what TestExtractDefeatsNoDefeat — tests extraction when no defeat edges exist
package analyzer

import (
	"testing"
)

func TestExtractDefeatsNoDefeat(t *testing.T) {
	path := writeGoFile(t, `package example

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func W(c, g any) (bool, any) { return true, nil }

var g = toulmin.NewGraph("simple").
	Warrant(W, 1.0)
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(graphs))
	}
	if len(graphs[0].Defeats) != 0 {
		t.Errorf("expected 0 defeats, got %d", len(graphs[0].Defeats))
	}
}
