//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseEndpointItem_Run — tests parseEndpointItem for the run/check child error, unrecognized statement, and run/check success branches
package parser

import (
	"strings"
	"testing"
)

func TestParseEndpointItem_Run(t *testing.T) {
	t.Run("RunCheckError", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 9,
			Children: []item{
				{Text: "run notbacktick", Line: 10},
			},
		}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseRunCheckItem")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted case name after") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Unrecognized", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 11,
			Children: []item{
				{Text: "unrecognized child statement", Line: 12},
			},
		}
		_, err := parseEndpointItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for an unrecognized statement")
		}
		if !strings.Contains(err.Error(), "unrecognized provides statement") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("RunSuccess", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 13,
			Children: []item{
				{Text: "run `case1`", Line: 14},
			},
		}
		ep, err := parseEndpointItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ep.Runs) != 1 || ep.Runs[0] != "case1" {
			t.Fatalf("expected Runs=[case1], got %+v", ep.Runs)
		}
		if len(ep.Checks) != 0 {
			t.Fatalf("expected empty Checks, got %+v", ep.Checks)
		}
	})

	t.Run("CheckSuccess", func(t *testing.T) {
		it := item{
			Text: "provides `x`",
			Line: 15,
			Children: []item{
				{Text: "check `case2`", Line: 16},
			},
		}
		ep, err := parseEndpointItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ep.Checks) != 1 || ep.Checks[0] != "case2" {
			t.Fatalf("expected Checks=[case2], got %+v", ep.Checks)
		}
		if len(ep.Runs) != 0 {
			t.Fatalf("expected empty Runs, got %+v", ep.Runs)
		}
	})
}
