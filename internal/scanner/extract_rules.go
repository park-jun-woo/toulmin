//ff:func feature=scanner type=scanner control=iteration dimension=1
//ff:what ExtractRules — parses a Go file and extracts rule declarations via AST
package scanner

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// ExtractRules parses a Go source file and returns the package name
// and rule declarations (func name + //rule: annotation lines).
func ExtractRules(path string) (string, []RuleDecl, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}
	var decls []RuleDecl
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Doc == nil {
			continue
		}
		lines := extractRuleLines(fn.Doc)
		if len(lines) == 0 {
			continue
		}
		decls = append(decls, RuleDecl{FuncName: fn.Name.Name, Lines: lines})
	}
	return f.Name.Name, decls, nil
}
