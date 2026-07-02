//ff:func feature=tangl type=parser control=sequence
//ff:what parseConstValue — split a ConstDef's literal value from its "as Ref" clause
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseConstValue splits "<literal> [as <Ref>]" into the raw literal text
// and an optional spec Ref, e.g. "650 as `credit`.`Threshold`".
func parseConstValue(s string, path string, line int) (string, *ast.Ref, error) {
	idx := strings.Index(s, " as `")
	if idx < 0 {
		return strings.TrimSpace(s), nil, nil
	}
	value := strings.TrimSpace(s[:idx])
	refPart := strings.TrimSpace(s[idx+len(" as "):])
	ref, rest, ok := parseRef(refPart)
	if !ok {
		return "", nil, errAt(path, line, "expected reference after 'as'")
	}
	if strings.TrimSpace(rest) != "" {
		return "", nil, errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return value, &ref, nil
}
