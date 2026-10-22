package parser

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func parseFileDefinition(name, rest, path string, line int) (ast.Definition, bool, error) {
	if !strings.HasPrefix(strings.TrimSpace(rest), "reads ") {
		return ast.Definition{}, false, nil
	}
	rest = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(rest), "reads"))
	typeName, rest, ok := takeBacktick(rest)
	if !ok || typeName == "" {
		return ast.Definition{}, true, errAt(path, line, "expected backtick-quoted loader type after 'reads'")
	}
	rest, ok = takeKeyword(rest, "from")
	if !ok {
		return ast.Definition{}, true, errAt(path, line, "expected 'from' after loader type")
	}
	rest = strings.TrimSpace(rest)
	filePath, err := strconv.Unquote(rest)
	if err != nil || filePath == "" {
		return ast.Definition{}, true, errAt(path, line, "expected one non-empty quoted file path")
	}
	return ast.Definition{Name: name, Kind: ast.FileSpecDef, LoadType: typeName, FilePath: filePath, Line: line}, true, nil
}
