package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func writeGoFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestExtractDefeatsMethonChain(t *testing.T) {
	path := writeGoFile(t, `package example

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func W(c, g any) (bool, any) { return true, nil }
func R(c, g any) (bool, any) { return true, nil }

var g = toulmin.NewGraph("check").
	Warrant(W, 1.0).
	Rebuttal(R, 1.0).
	Defeat(R, W)
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(graphs))
	}
	dg := graphs[0]
	if dg.Name != "check" {
		t.Errorf("expected graph name 'check', got %q", dg.Name)
	}
	if len(dg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(dg.Rules))
	}
	attackers, ok := dg.Defeats["W"]
	if !ok || len(attackers) != 1 || attackers[0] != "R" {
		t.Errorf("expected Defeats[W]=[R], got %v", dg.Defeats)
	}
}

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
		names[g.Name] = true
	}
	if !names["first"] || !names["second"] {
		t.Errorf("expected graphs 'first' and 'second', got %v", names)
	}
}

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
