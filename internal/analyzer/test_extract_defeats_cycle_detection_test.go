//ff:func feature=analyzer type=analyzer control=sequence
//ff:what TestExtractDefeatsCycleDetection — tests that cyclic defeat graphs are detected
package analyzer

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestExtractDefeatsCycleDetection(t *testing.T) {
	path := writeGoFile(t, `package example

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func A(c, g any) (bool, any) { return true, nil }
func B(c, g any) (bool, any) { return true, nil }

var g = toulmin.NewGraph("cyclic").
	Warrant(A, 1.0).
	Rebuttal(B, 1.0).
	Defeat(B, A).
	Defeat(A, B)
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(graphs))
	}
	if err := toulmin.DetectCycle(graphs[0].Defeats); err == nil {
		t.Fatal("expected cycle detection error")
	}
}
