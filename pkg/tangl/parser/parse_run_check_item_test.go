//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseRunCheckItem — tests parseRunCheckItem for run/check/unmatched branches
package parser

import "testing"

func TestParseRunCheckItem(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		it := item{Text: "run `case1`", Line: 1}
		name, kind, ok, err := parseRunCheckItem(it, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "case1" || kind != "run" {
			t.Errorf("expected name=case1 kind=run, got name=%q kind=%q", name, kind)
		}
	})

	t.Run("Check", func(t *testing.T) {
		it := item{Text: "check `case2`", Line: 2}
		name, kind, ok, err := parseRunCheckItem(it, "test.md")
		if !ok {
			t.Fatal("expected ok=true")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "case2" || kind != "check" {
			t.Errorf("expected name=case2 kind=check, got name=%q kind=%q", name, kind)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		it := item{Text: "something else", Line: 3}
		name, kind, ok, err := parseRunCheckItem(it, "test.md")
		if ok {
			t.Fatal("expected ok=false for an unrecognized statement")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "" || kind != "" {
			t.Errorf("expected empty name/kind, got name=%q kind=%q", name, kind)
		}
	})
}
