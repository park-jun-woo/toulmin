package ko

import (
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

func parseFileDefinitionKo(name, rest, path string, line int) (ast.Definition, bool, error) {
	if !strings.HasSuffix(strings.TrimSpace(rest), "파일에서 읽는다") {
		return ast.Definition{}, false, nil
	}
	body := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(rest), "파일에서 읽는다"))
	typeName, tail, ok := takeBacktick(body)
	if !ok || typeName == "" {
		return ast.Definition{}, true, errAtKo(path, line, "expected backtick-quoted loader type")
	}
	tail = skipParticlesKo(tail)
	tail = strings.TrimSpace(tail)
	filePath, err := strconv.Unquote(tail)
	if err != nil || filePath == "" {
		return ast.Definition{}, true, errAtKo(path, line, "expected one non-empty quoted file path")
	}
	return ast.Definition{Name: name, Kind: ast.FileSpecDef, LoadType: typeName, FilePath: filePath, Line: line}, true, nil
}
