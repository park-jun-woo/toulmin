//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEval — test eval declaration parsing
package parser

import "testing"

// TestParseEval tests eval declaration parsing with options.
func TestParseEval(t *testing.T) {
	ed, err := parseEval("acResult is results of evaluating access with trace", 20)
	if err != nil {
		t.Fatalf("parseEval failed: %v", err)
	}
	if ed.Name != "acResult" {
		t.Errorf("expected name 'acResult', got %q", ed.Name)
	}
	if ed.Graph != "access" {
		t.Errorf("expected graph 'access', got %q", ed.Graph)
	}
	if !ed.Trace {
		t.Error("expected trace=true")
	}
	if ed.Duration {
		t.Error("expected duration=false")
	}

	ed2, err := parseEval("result is results of evaluating myGraph", 21)
	if err != nil {
		t.Fatalf("parseEval without option failed: %v", err)
	}
	if ed2.Graph != "myGraph" {
		t.Errorf("expected graph 'myGraph', got %q", ed2.Graph)
	}
	if ed2.Trace || ed2.Duration {
		t.Error("expected no options")
	}

	ed3, err := parseEval("r is results of evaluating g with duration", 22)
	if err != nil {
		t.Fatalf("parseEval with duration failed: %v", err)
	}
	if !ed3.Duration {
		t.Error("expected duration=true")
	}
}
