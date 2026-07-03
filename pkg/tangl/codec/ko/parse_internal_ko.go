//ff:func feature=tangl type=parser control=sequence
//ff:what parseInternalHeadKo — parse an internal trigger head (발생 시 event / 매 ~마다 tick)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseInternalHeadKo parses the tangl:내부 trigger head: "`이벤트` 발생 시" (OnEvent,
// the backtick name is the event) or "매 <interval>마다[, `C` 까지]" (EveryTick,
// parseEveryClauseKo normalizes the interval and reads the until case). It returns
// an ast.Internal with no body — the caller walks the run/check children.
func parseInternalHeadKo(it scan.Item, path string) (ast.Internal, error) {
	if body, ok := takeSuffixKo(it.Text, "발생 시"); ok {
		name, rest, ok := takeBacktick(body)
		if !ok {
			return ast.Internal{}, errAtKo(path, it.No, "expected backtick-quoted event before '발생 시'")
		}
		if strings.TrimSpace(skipParticlesKo(rest)) != "" {
			return ast.Internal{}, errAtKo(path, it.No, "unexpected trailing text: %q", rest)
		}
		return ast.Internal{Kind: ast.OnEvent, Event: name, Line: it.No}, nil
	}
	rest, ok := takeKeywordKo(it.Text, "매")
	if !ok {
		return ast.Internal{}, errAtKo(path, it.No, "expected '`이벤트` 발생 시' or '매 <간격>마다', got %q", it.Text)
	}
	interval, until, err := parseEveryClauseKo(rest, it.No, path)
	if err != nil {
		return ast.Internal{}, err
	}
	return ast.Internal{Kind: ast.EveryTick, Interval: interval, Until: until, Line: it.No}, nil
}
