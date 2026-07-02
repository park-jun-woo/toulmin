//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseInternalHead — tests parseInternalHead for on/every/unmatched branches
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func TestParseInternalHead(t *testing.T) {
	t.Run("OnEmptyEvent", func(t *testing.T) {
		it := item{Text: "on   ", Line: 1}
		_, err := parseInternalHead(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for an empty event")
		}
		if !strings.Contains(err.Error(), "expected event after 'on'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("OnSuccess", func(t *testing.T) {
		it := item{Text: "on start", Line: 2}
		internal, err := parseInternalHead(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if internal.Kind != ast.OnEvent {
			t.Errorf("expected OnEvent, got %v", internal.Kind)
		}
		if internal.Event != "start" {
			t.Errorf("expected Event=start, got %q", internal.Event)
		}
	})

	t.Run("EveryError", func(t *testing.T) {
		it := item{Text: "every   ", Line: 3}
		_, err := parseInternalHead(it, "test.md")
		if err == nil {
			t.Fatal("expected an error from parseEveryClause")
		}
		if !strings.Contains(err.Error(), "expected interval after 'every'") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("EverySuccess", func(t *testing.T) {
		it := item{Text: "every 1 day until `done`", Line: 4}
		internal, err := parseInternalHead(it, "test.md")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if internal.Kind != ast.EveryTick {
			t.Errorf("expected EveryTick, got %v", internal.Kind)
		}
		if internal.Interval != "1 day" {
			t.Errorf("expected Interval='1 day', got %q", internal.Interval)
		}
		if internal.Until != "done" {
			t.Errorf("expected Until=done, got %q", internal.Until)
		}
	})

	t.Run("Unmatched", func(t *testing.T) {
		it := item{Text: "when something happens", Line: 5}
		_, err := parseInternalHead(it, "test.md")
		if err == nil {
			t.Fatal("expected an error for an unrecognized head")
		}
		if !strings.Contains(err.Error(), "expected 'on <event>' or 'every <interval>'") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
