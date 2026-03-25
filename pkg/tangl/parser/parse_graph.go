//ff:func feature=tangl type=parser control=sequence
//ff:what parseGraph — parse a graph declaration into GraphDecl AST node
package parser

import (
	"fmt"
	"strings"
)

// parseGraph parses: name is a graph "id"
func parseGraph(text string, lineNum int) (GraphDecl, error) {
	parts := strings.SplitN(text, " is a graph ", 2)
	if len(parts) != 2 {
		return GraphDecl{}, fmt.Errorf("invalid graph declaration: %s", text)
	}
	name := strings.TrimSpace(parts[0])
	id := strings.Trim(strings.TrimSpace(parts[1]), "\"")
	if name == "" || id == "" {
		return GraphDecl{}, fmt.Errorf("invalid graph declaration: empty name or id in %s", text)
	}
	return GraphDecl{Name: name, ID: id, Line: lineNum}, nil
}
