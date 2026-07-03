//ff:func feature=tangl type=parser control=sequence
//ff:what parseNodeko — parse '`X` 는 [<fn> 쓰는] (일반|반박|예외) 규칙이다[꼬리]' into an ast.Node
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseNodeKo parses a case node registration. 규칙이다 sits mid-statement (기준/가중
// tails may follow), so it is located by contains-keyword dispatch; the role word
// (일반/반박/예외) is the last token of the head before 규칙이다, and any '<fn> 쓰는'
// clause precedes it (a function ref via applyNodeUsingKo, or '`Y` 판정' checking).
// The trailing '기준'/'가중' clauses are parsed by applyNodeTailKo. It returns
// ok=false (no error) when the item is not a node declaration.
func parseNodeKo(it scan.Item, path string) (ast.Node, bool, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Node{}, false, nil
	}
	rest = skipParticlesKo(rest)
	idx := strings.Index(rest, "규칙이다")
	if idx < 0 {
		return ast.Node{}, false, nil
	}
	head := strings.TrimSpace(rest[:idx])
	tail := strings.TrimSpace(rest[idx+len("규칙이다"):])
	roleWord := lastTokenKo(head)
	role, ok := parseRoleKo(roleWord)
	if !ok {
		return ast.Node{}, true, errAtKo(path, it.No, "expected role (일반/반박/예외) before '규칙이다', got %q", head)
	}
	before := strings.TrimSpace(strings.TrimSuffix(head, roleWord))
	n := ast.Node{Name: name, Role: role, Line: it.No}
	if before != "" {
		usingPart, ok := takeSuffixKo(before, "쓰는")
		if !ok {
			return ast.Node{}, true, errAtKo(path, it.No, "expected '쓰는' before the role, got %q", before)
		}
		if err := applyNodeUsingKo(&n, usingPart, it.No, path); err != nil {
			return ast.Node{}, true, err
		}
	}
	if tail != "" {
		if err := applyNodeTailKo(&n, tail, it.No, path); err != nil {
			return ast.Node{}, true, err
		}
	}
	return n, true, nil
}
