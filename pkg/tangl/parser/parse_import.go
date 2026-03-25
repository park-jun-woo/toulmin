//ff:func feature=tangl type=parser control=sequence
//ff:what parseImport — parse an import statement into Import AST node
package parser

import (
	"fmt"
	"strings"
)

// parseImport parses: alias is from "package/path"
func parseImport(text string, lineNum int) (Import, error) {
	parts := strings.SplitN(text, " is from ", 2)
	if len(parts) != 2 {
		return Import{}, fmt.Errorf("invalid import: %s", text)
	}
	alias := strings.TrimSpace(parts[0])
	pkg := strings.TrimSpace(parts[1])
	pkg = strings.Trim(pkg, "\"")
	if alias == "" || pkg == "" {
		return Import{}, fmt.Errorf("invalid import: empty alias or package in %s", text)
	}
	return Import{Alias: alias, Package: pkg, Line: lineNum}, nil
}
