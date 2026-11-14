//ff:func feature=tangl type=parser control=sequence
//ff:what parseHostRuleKo — parse a 규칙 host declaration ('`n` 은 ... `별칭`.`함수` 쓰는 규칙이다')
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseHostRuleKo parses a host-function rule declaration ending in 규칙이다. The
// bool reports whether the entry is a host declaration at all (false leaves it to
// the condition-rule parser). The declaration carries the host function (쓰는),
// optional criteria (기준으로) and grounds (근거로); weight and level are rejected
// because they belong to a case binding.
func parseHostRuleKo(it scan.Item, path string) (ast.InlineRule, bool, error) {
	body, ok := takeSuffixKo(it.Text, "규칙이다")
	if !ok {
		return ast.InlineRule{}, false, nil
	}
	name, rest, ok := takeBacktick(body)
	if !ok {
		return ast.InlineRule{}, true, errAtKo(path, it.No, "expected backtick-quoted rule name")
	}
	n := ast.Node{}
	rest, err := parseNodeModifiersKo(&n, rest, it.No, path)
	if err != nil {
		return ast.InlineRule{}, true, err
	}
	if grounds, tail, ok := parseGroundModifierKo(rest); ok {
		n.Grounds = grounds
		rest = strings.TrimSpace(skipParticlesKo(tail))
	}
	usingPart, ok := takeSuffixKo(rest, "쓰는")
	if !ok {
		return ast.InlineRule{}, true, errAtKo(path, it.No, "expected function reference followed by '쓰는'")
	}
	if err := applyNodeUsingKo(&n, usingPart, it.No, path); err != nil {
		return ast.InlineRule{}, true, err
	}
	if n.Using == nil || n.Checking != nil {
		return ast.InlineRule{}, true, errAtKo(path, it.No, "rule declaration requires a host function, not a checking case")
	}
	if n.Qualified != nil || n.Level != ast.LevelUnset {
		return ast.InlineRule{}, true, errAtKo(path, it.No, "weight and level belong to a case binding")
	}
	return ast.InlineRule{Name: name, Using: n.Using, With: n.With, Grounds: n.Grounds, Line: it.No}, true, nil
}
