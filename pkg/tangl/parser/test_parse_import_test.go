//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseImport — test import statement parsing
package parser

import "testing"

// TestParseImport tests import parsing.
func TestParseImport(t *testing.T) {
	imp, err := parseImport(`policy is from "github.com/example/pkg"`, 1)
	if err != nil {
		t.Fatalf("parseImport failed: %v", err)
	}
	if imp.Alias != "policy" {
		t.Errorf("expected alias 'policy', got %q", imp.Alias)
	}
	if imp.Package != "github.com/example/pkg" {
		t.Errorf("expected package 'github.com/example/pkg', got %q", imp.Package)
	}
	if imp.Line != 1 {
		t.Errorf("expected line 1, got %d", imp.Line)
	}
}
