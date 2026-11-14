//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEndpointItem — tests parseEndpointItem for header parsing errors and the no-children/require branches
package parser

import (
	"strings"
	"testing"
)

func TestParseEndpointItem(t *testing.T) {
	t.Run("NoProvides", func(t *testing.T) {
		it := item{Text: "gives `x`", Line: 1}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing 'provides'")
		}
		if !strings.Contains(err.Error(), "expected 'provides `name`'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoBacktickName", func(t *testing.T) {
		it := item{Text: "provides x", Line: 2}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for missing backtick-quoted endpoint name")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted endpoint name") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("TrailingText", func(t *testing.T) {
		it := item{Text: "provides `x` extra", Line: 3}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for trailing text")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoChildren", func(t *testing.T) {
		it := item{Text: "provides `x`", Line: 4}
		ep, err := parseEndpointItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ep.Name != "x" {
			t.Errorf("expected Name=x, got %q", ep.Name)
		}
		if len(ep.Requires) != 0 || len(ep.Runs) != 0 || len(ep.Checks) != 0 {
			t.Errorf("expected empty endpoint children, got %+v", ep)
		}
	})

	t.Run("RequireError", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 5,
			Children: []item{
				{Text: "`f` is required extra", Line: 6},
			},
		}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseRequireItem")
		}
		if !strings.Contains(err.Error(), "unexpected trailing text") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("RequireSuccess", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 7,
			Children: []item{
				{Text: "`f` is required", Line: 8},
			},
		}
		ep, err := parseEndpointItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ep.Requires) != 1 || ep.Requires[0].Field != "f" {
			t.Fatalf("expected one require for field f, got %+v", ep.Requires)
		}
	})
}
