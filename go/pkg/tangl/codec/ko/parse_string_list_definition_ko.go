package ko

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func parseStringListDefinitionKo(name, source, path string, line int) (ast.Definition, bool, error) {
	if !strings.HasPrefix(strings.TrimSpace(source), "[") {
		return ast.Definition{}, false, nil
	}
	var values []string
	dec := json.NewDecoder(strings.NewReader(source))
	if err := dec.Decode(&values); err != nil {
		return ast.Definition{}, true, errAtKo(path, line, "invalid JSON string array: %v", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return ast.Definition{}, true, errAtKo(path, line, "unexpected value after JSON string array")
	}
	return ast.Definition{Name: name, Kind: ast.StringListDef, Strings: values, Line: line}, true, nil
}
