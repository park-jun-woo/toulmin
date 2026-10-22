//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInternalItem — tests parseInternalItem for head errors and children run/check branches
package parser

import (
	"strings"
	"testing"
)

func TestParseInternalItem(t *testing.T) {
	t.Run("HeadError", func(t *testing.T) {
		it := item{Text: "bogus head", Line: 1}
		_, err := parseInternalItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseInternalHead")
		}
		if !strings.Contains(err.Error(), "expected 'on <event>' or 'every <interval>'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("NoChildren", func(t *testing.T) {
		it := item{Text: "on start", Line: 2}
		in, err := parseInternalItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(in.Runs) != 0 || len(in.Checks) != 0 {
			t.Errorf("expected empty runs/checks, got %+v", in)
		}
	})

	t.Run("ChildError", func(t *testing.T) {
		it := item{
			Text: "on start",
			Line: 3,
			Children: []item{
				{Text: "run notbacktick", Line: 4},
			},
		}
		_, err := parseInternalItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseRunCheckItem")
		}
		if !strings.Contains(err.Error(), "expected backtick-quoted case name after") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ChildUnrecognized", func(t *testing.T) {
		it := item{
			Text: "on start",
			Line: 5,
			Children: []item{
				{Text: "unrecognized child statement", Line: 6},
			},
		}
		_, err := parseInternalItem(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for an unrecognized statement")
		}
		if !strings.Contains(err.Error(), "unrecognized internal statement") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ChildRun", func(t *testing.T) {
		it := item{
			Text: "on start",
			Line: 7,
			Children: []item{
				{Text: "run `case1`", Line: 8},
			},
		}
		in, err := parseInternalItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(in.Runs) != 1 || in.Runs[0] != "case1" {
			t.Fatalf("expected Runs=[case1], got %+v", in.Runs)
		}
		if len(in.Checks) != 0 {
			t.Fatalf("expected empty Checks, got %+v", in.Checks)
		}
	})

	t.Run("ChildCheck", func(t *testing.T) {
		it := item{
			Text: "on start",
			Line: 9,
			Children: []item{
				{Text: "check `case2`", Line: 10},
			},
		}
		in, err := parseInternalItem(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(in.Checks) != 1 || in.Checks[0] != "case2" {
			t.Fatalf("expected Checks=[case2], got %+v", in.Checks)
		}
		if len(in.Runs) != 0 {
			t.Fatalf("expected empty Runs, got %+v", in.Runs)
		}
	})
}
