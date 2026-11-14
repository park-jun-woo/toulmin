//ff:func feature=tangl type=parser control=sequence
//ff:what parseStringListDefinition — parse a JSON string-array definition body into a StringListDef
package parser

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func parseStringListDefinition(name, source, path string, line int) (ast.Definition, bool, error) {
	if !strings.HasPrefix(strings.TrimSpace(source), "[") {
		return ast.Definition{}, false, nil
	}
	var values []string
	dec := json.NewDecoder(strings.NewReader(source))
	if err := dec.Decode(&values); err != nil {
		return ast.Definition{}, true, errAt(path, line, "invalid JSON string array: %v", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return ast.Definition{}, true, errAt(path, line, "unexpected value after JSON string array")
	}
	return ast.Definition{Name: name, Kind: ast.StringListDef, Strings: values, Line: line}, true, nil
}
