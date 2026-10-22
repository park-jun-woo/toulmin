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
	if def, matched, err := parseFileDefinition(name, rest, path, it.Line); matched {
		if len(it.Children) != 0 && err == nil {
			return ast.Definition{}, errAt(path, it.Line, "file definition cannot have child fields")
		}
		return def, err
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
	if def, matched, err := parseStringListDefinition(name, rest, path, it.Line); matched {
		if len(it.Children) != 0 && err == nil {
			return ast.Definition{}, errAt(path, it.Line, "string-list definition cannot have child fields")
		}
		return def, err
	}
	value, specRef, err := parseConstValue(rest, path, it.Line)
	if err != nil {
		return ast.Definition{}, err
	}
	return ast.Definition{Name: name, Kind: ast.ConstDef, Value: value, SpecRef: specRef, Line: it.Line}, nil
}
