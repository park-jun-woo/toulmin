//ff:func feature=tangl type=parser control=sequence
//ff:what TestParseNodeItem — tests parseNodeItem for name/is/article/role no-match branches
package parser

import (
	"strings"
	"testing"
)

func TestParseNodeItem(t *testing.T) {
	t.Run("NoBacktickName", func(t *testing.T) {
		it := item{Text: "node1 is a general rule", Line: 1}
		n, matched, err := parseNodeItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing backtick-quoted name")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Name != "" || n.Line != 0 {
			t.Errorf("expected zero-value Node, got %+v", n)
		}
	})

	t.Run("NoIs", func(t *testing.T) {
		it := item{Text: "`node1` becomes a general rule", Line: 2}
		n, matched, err := parseNodeItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing 'is'")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Name != "" || n.Line != 0 {
			t.Errorf("expected zero-value Node, got %+v", n)
		}
	})

	t.Run("NoArticle", func(t *testing.T) {
		it := item{Text: "`node1` is general rule", Line: 3}
		n, matched, err := parseNodeItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for missing article")
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n.Name != "" || n.Line != 0 {
			t.Errorf("expected zero-value Node, got %+v", n)
		}
	})

	t.Run("BadRole", func(t *testing.T) {
		it := item{Text: "`node1` is a strange rule", Line: 4}
		n, matched, err := parseNodeItem(it, "test.md")
		if matched {
			t.Fatal("expected matched=false for a bad role")
		}
		if err == nil {
			t.Fatal("expected an error for a bad role")
		}
		if !strings.Contains(err.Error(), "expected role (general/counter/except rule)") {
			t.Errorf("unexpected error: %v", err)
		}
		if n.Name != "" || n.Line != 0 {
			t.Errorf("expected zero-value Node, got %+v", n)
		}
	})
}
