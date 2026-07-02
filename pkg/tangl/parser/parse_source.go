//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what ParseSource — parse TANGL v0.3 markdown source into a Document
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// ParseSource parses TANGL v0.3 markdown source (src) into an ast.Document.
// path is used only for error messages ("path:line: message"). Parsing
// stops and returns the first error encountered (no error collection).
func ParseSource(src, path string) (*ast.Document, error) {
	lines := splitLines(src)
	secs := splitSections(lines)
	doc := &ast.Document{Path: path}
	maxOrder := -1
	for _, sec := range secs {
		idx, ok := sectionOrderIndex(sec.Name)
		if !ok {
			return nil, errAt(path, sec.HeaderLine, "unknown tangl section: %q", sec.Name)
		}
		if idx < maxOrder {
			return nil, errAt(path, sec.HeaderLine, "tangl:%s appears out of order", sec.Name)
		}
		maxOrder = idx
		if err := applySection(doc, sec, path); err != nil {
			return nil, err
		}
	}
	if doc.Subject == "" {
		return nil, errAt(path, 1, "tangl:Subject section is required")
	}
	if len(doc.Cases) == 0 {
		return nil, errAt(path, 1, "tangl:Cases section is required")
	}
	return doc, nil
}
