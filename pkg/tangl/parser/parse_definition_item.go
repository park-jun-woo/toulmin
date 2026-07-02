//ff:func feature=tangl type=parser control=sequence
//ff:what parseDefinitionItem — parse one Definitions entry (const or struct)
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseDefinitionItem parses "`x` means <literal> [as <Ref>]" (ConstDef) or
// "`x` means" followed by nested "has `f` as Type" children (StructDef).
func parseDefinitionItem(it item, path string) (ast.Definition, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Definition{}, errAt(path, it.Line, "expected backtick-quoted term name, got %q", it.Text)
	}
	rest, ok = takeKeyword(rest, "means")
	if !ok {
		return ast.Definition{}, errAt(path, it.Line, "expected 'means' after term name")
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		fields, err := parseDefinitionFields(it.Children, path)
		if err != nil {
			return ast.Definition{}, err
		}
		return ast.Definition{Name: name, Kind: ast.StructDef, Fields: fields, Line: it.Line}, nil
	}
	value, specRef, err := parseConstValue(rest, path, it.Line)
	if err != nil {
		return ast.Definition{}, err
	}
	return ast.Definition{Name: name, Kind: ast.ConstDef, Value: value, SpecRef: specRef, Line: it.Line}, nil
}
