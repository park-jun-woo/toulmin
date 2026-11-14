//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseNodeItemLevel — a role followed by ", level X" sets the level; bad levels and trailing text are errors
package parser

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestParseNodeItemLevel(t *testing.T) {
	n, ok, err := parseNodeItem(item{Text: "`x` is a counter rule, level fail", Line: 5}, "t.md")
	if err != nil || !ok || n.Name != "x" || n.Role != ast.CounterRule || n.Level == ast.LevelUnset || n.Line != 5 {
		t.Fatalf("level clause: %+v ok=%v err=%v", n, ok, err)
	}
	for _, text := range []string{
		"`x` is a counter rule, level nonsense",
		"`x` is a counter rule, level fail extra",
	} {
		_, ok, err := parseNodeItem(item{Text: text, Line: 5}, "t.md")
		if !ok || err == nil {
			t.Errorf("%q: ok=%v err=%v, want a matched error", text, ok, err)
		}
	}
	_, _, err = parseNodeItem(item{Text: "`x` is a counter rule, level fail extra", Line: 5}, "t.md")
	if err == nil || !strings.Contains(err.Error(), "unexpected trailing text") {
		t.Fatalf("trailing text: %v", err)
	}
}
