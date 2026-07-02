//ff:func feature=tangl type=parser control=sequence
//ff:what Parse — read a TANGL v0.3 markdown file and parse it into a Document
package parser

import (
	"os"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// Parse reads the file at path and parses it into an ast.Document.
func Parse(path string) (*ast.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseSource(string(data), path)
}
