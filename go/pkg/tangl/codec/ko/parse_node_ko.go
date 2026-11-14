//ff:func feature=tangl type=parser control=sequence
//ff:what parseNodeKo — parse '`X` 은 [관형절들] [<fn> 쓰는] (일반|반박|예외) 규칙이다' into an ast.Node
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseNodeKo parses a case node registration. With the tail clauses gone, 규칙이다
// is the sentence-final keyword, so the statement is located by suffix dispatch
// (takeSuffixKo) like every other KO statement. The role word (일반/반박/예외) is the
// last token of the head; the adnominal modifiers before it — '`spec` 을 기준으로 하는'
// (With), '가중치 <float> 의' (Qualified), '(실패를|검토를) 보고하는' (Level) — are consumed
// by parseNodeModifiersKo, and any remaining '<fn> 쓰는' clause by applyNodeUsingKo
// (a function ref, or '`Y` 판정' checking). It returns ok=false (no error) when the
// item is not a node declaration.
func parseNodeKo(it scan.Item, path string) (ast.Node, bool, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Node{}, false, nil
	}
	head, ok := takeSuffixKo(skipParticlesKo(rest), "규칙이다")
	if !ok {
		return ast.Node{}, false, nil
	}
	roleWord := lastTokenKo(head)
	role, ok := parseRoleKo(roleWord)
	if !ok {
		return ast.Node{}, true, errAtKo(path, it.No, "expected role (일반/반박/예외) before '규칙이다', got %q", head)
	}
	before := strings.TrimSpace(strings.TrimSuffix(head, roleWord))
	n := ast.Node{Name: name, Role: role, Line: it.No}
	rest, err := parseNodeModifiersKo(&n, before, it.No, path)
	if err != nil {
		return ast.Node{}, true, err
	}
	if rest != "" {
		usingPart, ok := takeSuffixKo(rest, "쓰는")
		if !ok {
			return ast.Node{}, true, errAtKo(path, it.No, "expected '쓰는' before the role, got %q", rest)
		}
		if err := applyNodeUsingKo(&n, usingPart, it.No, path); err != nil {
			return ast.Node{}, true, err
		}
	}
	return n, true, nil
}
