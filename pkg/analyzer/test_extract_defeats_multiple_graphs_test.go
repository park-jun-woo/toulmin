//ff:func feature=analyzer type=analyzer control=iteration dimension=1
//ff:what TestExtractDefeatsMultipleGraphs — tests extraction of multiple Graph builders from one file
package analyzer

import (
	"testing"
)

func TestExtractDefeatsMultipleGraphs(t *testing.T) {
	path := writeGoFile(t, `package example

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func W(c, g any) (bool, any) { return true, nil }
func R(c, g any) (bool, any) { return true, nil }
func D(c, g any) (bool, any) { return true, nil }

var g1 = toulmin.NewGraph("first").
	Warrant(W, 1.0).
	Rebuttal(R, 1.0).
	Defeat(R, W)

var g2 = toulmin.NewGraph("second").
	Warrant(W, 1.0).
	Defeater(D, 1.0).
	Defeat(D, W)
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 2 {
		t.Fatalf("expected 2 graphs, got %d", len(graphs))
	}
	names := map[string]bool{}
	for _, g := range graphs {
		names[g.Graph] = true
	}
	if !names["first"] || !names["second"] {
		t.Errorf("expected graphs 'first' and 'second', got %v", names)
	}
}
