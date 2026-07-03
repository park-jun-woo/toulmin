//ff:func feature=tangl type=parser control=sequence
//ff:what applyNodeUsingKo — resolve a node's '<fn> 쓰는' clause to Using ref or Checking case
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// applyNodeUsingKo parses the clause before '쓰는' on a node declaration: either
// '`Y` 판정' (verdict composition → n.Checking) or a function reference
// '`fn`' / '`alias`.`fn`' (→ n.Using). It mirrors pkg/tangl/parser.applyNodeClause's
// using/checking split.
func applyNodeUsingKo(n *ast.Node, usingPart string, line int, path string) error {
	if before, ok := takeSuffixKo(usingPart, "판정"); ok {
		name, rest, ok := takeBacktick(before)
		if !ok {
			return errAtKo(path, line, "expected backtick-quoted case name before '판정'")
		}
		if strings.TrimSpace(skipParticlesKo(rest)) != "" {
			return errAtKo(path, line, "unexpected text before '판정': %q", before)
		}
		n.Checking = name
		return nil
	}
	ref, rest, ok := parseRefKo(usingPart)
	if !ok {
		return errAtKo(path, line, "expected function reference before '쓰는', got %q", usingPart)
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return errAtKo(path, line, "unexpected trailing text after function reference: %q", rest)
	}
	n.Using = &ref
	return nil
}
