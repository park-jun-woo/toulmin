//ff:func feature=tangl type=parser control=sequence
//ff:what parseInternalHead — parse an Internal item's "on"/"every" trigger head
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseInternalHead parses the trigger portion of an Internal item, before
// its nested run/check children: "on <event>" or "every <interval> [until `case`]".
func parseInternalHead(it item, path string) (ast.Internal, error) {
	if rest, ok := takeKeyword(it.Text, "on"); ok {
		event := strings.TrimSpace(rest)
		if event == "" {
			return ast.Internal{}, errAt(path, it.Line, "expected event after 'on'")
		}
		return ast.Internal{Kind: ast.OnEvent, Event: event, Line: it.Line}, nil
	}
	if rest, ok := takeKeyword(it.Text, "every"); ok {
		interval, until, err := parseEveryClause(rest, path, it.Line)
		if err != nil {
			return ast.Internal{}, err
		}
		return ast.Internal{Kind: ast.EveryTick, Interval: interval, Until: until, Line: it.Line}, nil
	}
	return ast.Internal{}, errAt(path, it.Line, "expected 'on <event>' or 'every <interval>', got %q", it.Text)
}
