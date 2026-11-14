//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseCaseItem — tests parseCaseItem for header parsing errors, empty/populated children, and child-error propagation
package parser

import (
	"strings"
	"testing"
)

func TestParseCaseItem(t *testing.T) {
	t.Run("MissingPrefix", func(t *testing.T) {
		_, err := parseCaseItem(item{Text: "not a case", Line: 1}, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected 'in case of") {
			t.Fatalf("expected prefix error, got %v", err)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		_, err := parseCaseItem(item{Text: "in case of notabacktick", Line: 2}, "test.md")
		if err == nil || !strings.Contains(err.Error(), "expected backtick-quoted case name") {
			t.Fatalf("expected name error, got %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		_, err := parseCaseItem(item{Text: "in case of `c1` extra", Line: 3}, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Fatalf("expected trailing text error, got %v", err)
		}
	})

	t.Run("NoChildren", func(t *testing.T) {
		c, err := parseCaseItem(item{Text: "in case of `c1`", Line: 4}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.Name != "c1" {
			t.Fatalf("expected Name c1, got %q", c.Name)
		}
		if len(c.Requires) != 0 || len(c.Nodes) != 0 || len(c.Attacks) != 0 || len(c.Execs) != 0 {
			t.Fatalf("expected no children applied, got %+v", c)
		}
	})

	t.Run("WithValidChild", func(t *testing.T) {
		c, err := parseCaseItem(item{
			Text: "in case of `c1`",
			Line: 5,
			Children: []item{
				{Text: "`amount` is required", Line: 6},
			},
		}, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(c.Requires) != 1 || c.Requires[0].Field != "amount" {
			t.Fatalf("expected Requires to contain amount, got %+v", c.Requires)
		}
	})

	t.Run("ChildErrorPropagates", func(t *testing.T) {
		_, err := parseCaseItem(item{
			Text: "in case of `c1`",
			Line: 7,
			Children: []item{
				{Text: "unrecognized child statement", Line: 8},
			},
		}, "test.md")
		if err == nil || !strings.Contains(err.Error(), "unrecognized case statement") {
			t.Fatalf("expected child error to propagate, got %v", err)
		}
	})
}
