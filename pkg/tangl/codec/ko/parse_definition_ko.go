//ff:func feature=tangl type=parser control=sequence
//ff:what parseDefinitionKo — parse one 정의 entry (const value or struct head)
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/codec/scan"
)

// parseDefinitionKo parses a tangl:정의 entry. 뜻한다 sits mid-statement (a
// ', 타입 <Ref>' tail may follow), so it is located by contains-keyword dispatch.
// "다음(을) 뜻한다" opens a StructDef whose fields are the item's children;
// otherwise the head is a ConstDef literal with an optional 타입 spec Ref.
func parseDefinitionKo(it scan.Item, path string) (ast.Definition, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Definition{}, errAtKo(path, it.No, "expected backtick-quoted term name, got %q", it.Text)
	}
	rest = skipParticlesKo(rest)
	idx := strings.Index(rest, "뜻한다")
	if idx < 0 {
		return ast.Definition{}, errAtKo(path, it.No, "expected '뜻한다' in definition, got %q", it.Text)
	}
	valuePart := strings.TrimSpace(rest[:idx])
	tail := strings.TrimSpace(rest[idx+len("뜻한다"):])
	if valuePart == "다음" || valuePart == "다음을" || valuePart == "다음를" {
		fields, err := parseDefinitionFieldsKo(it.Children, path)
		if err != nil {
			return ast.Definition{}, err
		}
		return ast.Definition{Name: name, Kind: ast.StructDef, Fields: fields, Line: it.No}, nil
	}
	value, specRef, err := parseConstValueKo(valuePart, tail, it.No, path)
	if err != nil {
		return ast.Definition{}, err
	}
	return ast.Definition{Name: name, Kind: ast.ConstDef, Value: value, SpecRef: specRef, Line: it.No}, nil
}
