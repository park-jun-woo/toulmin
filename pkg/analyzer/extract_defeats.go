//ff:func feature=analyzer type=analyzer control=iteration dimension=1
//ff:what ExtractDefeats — extracts defeat relationships from Go source via AST
package analyzer

import (
	"go/parser"
	"go/token"
)

// ExtractDefeats parses a Go source file and extracts GraphBuilder
// defeat relationships. Returns one DefeatGraph per NewGraph call found.
func ExtractDefeats(path string) ([]DefeatGraph, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	var graphs []DefeatGraph
	for _, call := range findGraphCalls(f) {
		name := extractGraphName(call)
		if name == "" {
			continue
		}
		dg := DefeatGraph{
			Name:    name,
			Defeats: make(map[string][]string),
		}
		collectChain(call, &dg)
		graphs = append(graphs, dg)
	}
	return graphs, nil
}
