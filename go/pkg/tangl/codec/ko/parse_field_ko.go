//ff:func feature=tangl type=parser control=sequence
//ff:what parseFieldKo — parse a StructDef field ('`f` 를 가진다, 타입 T')
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseFieldKo parses "`field` 를 가진다, 타입 <Type>". 가진다 is followed by a comma
// (not a space) before the 타입 tail, so it is located by contains-keyword
// dispatch and the tail is parsed after the comma is dropped as a particle.
func parseFieldKo(it scan.Item, path string) (ast.Field, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Field{}, errAtKo(path, it.No, "expected backtick-quoted field name, got %q", it.Text)
	}
	rest = skipParticlesKo(rest)
	idx := strings.Index(rest, "가진다")
	if idx < 0 {
		return ast.Field{}, errAtKo(path, it.No, "expected '가진다' in field, got %q", it.Text)
	}
	if strings.TrimSpace(rest[:idx]) != "" {
		return ast.Field{}, errAtKo(path, it.No, "unexpected text before '가진다': %q", rest[:idx])
	}
	tail := skipParticlesKo(strings.TrimSpace(rest[idx+len("가진다"):]))
	tail, ok = takeKeywordKo(tail, "타입")
	if !ok {
		return ast.Field{}, errAtKo(path, it.No, "expected '타입' after '가진다'")
	}
	typ, rest2, ok := parseTypeRefKo(strings.TrimSpace(tail))
	if !ok {
		return ast.Field{}, errAtKo(path, it.No, "expected type after '타입'")
	}
	if strings.TrimSpace(rest2) != "" {
		return ast.Field{}, errAtKo(path, it.No, "unexpected trailing text: %q", rest2)
	}
	return ast.Field{Name: name, Type: typ, Line: it.No}, nil
}
