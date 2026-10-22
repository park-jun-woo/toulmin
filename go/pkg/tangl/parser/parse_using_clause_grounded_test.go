//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseUsingClause_Grounded — tests parseUsingClause for grounded/combination/error branches
package parser

import (
	"strings"
	"testing"
)

func TestParseUsingClause_Grounded(t *testing.T) {
	t.Run("GroundedOnly", func(t *testing.T) {
		ref, with, qual, grounds, rest, err := parseUsingClause("`fn` grounded on `source-body`", "test.md", 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if with != nil || qual != nil {
			t.Errorf("expected nil with/qual, got %v / %v", with, qual)
		}
		if len(grounds) != 1 || grounds[0] != "source-body" {
			t.Errorf("expected grounds=[source-body], got %v", grounds)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})

	t.Run("WithAndGrounded", func(t *testing.T) {
		ref, with, _, grounds, rest, err := parseUsingClause("`fn` with `term1` grounded on `source-body` and `mx`", "test.md", 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if len(with) != 1 || with[0] != "term1" {
			t.Errorf("expected with=[term1], got %v", with)
		}
		if len(grounds) != 2 || grounds[0] != "source-body" || grounds[1] != "mx" {
			t.Errorf("expected grounds=[source-body mx], got %v", grounds)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})

	t.Run("WithQualifiedGrounded", func(t *testing.T) {
		ref, with, qual, grounds, rest, err := parseUsingClause("`fn` with `term1` qualified 0.5 grounded on `source-body`", "test.md", 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Name != "fn" {
			t.Errorf("expected ref.Name=fn, got %+v", ref)
		}
		if len(with) != 1 || with[0] != "term1" {
			t.Errorf("expected with=[term1], got %v", with)
		}
		if qual == nil || *qual != 0.5 {
			t.Errorf("expected qual=0.5, got %v", qual)
		}
		if len(grounds) != 1 || grounds[0] != "source-body" {
			t.Errorf("expected grounds=[source-body], got %v", grounds)
		}
		if rest != "" {
			t.Errorf("expected empty rest, got %q", rest)
		}
	})

	t.Run("GroundedError", func(t *testing.T) {
		_, _, _, _, _, err := parseUsingClause("`fn` grounded on notbacktick", "test.md", 4)
		if err == nil {
			t.Fatal("expected an error from parseGroundedClause")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted name after 'grounded on'/'and'") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
