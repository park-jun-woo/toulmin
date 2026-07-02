//ff:func feature=cli type=command control=sequence
//ff:what resolveEffectsTarget — resolves the effects command's file+endpoint or subject.endpoint+dir form
package tanglcli

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// resolveEffectsTarget resolves args into a parsed Document and an
// endpoint name. Two args ("<file.md> <endpoint>") parse the file
// directly; one arg ("<subject>.<endpoint>") scans dir for a document
// whose Subject matches.
func resolveEffectsTarget(args []string, dir string) (*ast.Document, string, error) {
	if len(args) == 2 {
		doc, err := parser.Parse(args[0])
		if err != nil {
			return nil, "", err
		}
		return doc, trimBackticks(args[1]), nil
	}
	subject, endpoint, err := splitSubjectEndpoint(args[0])
	if err != nil {
		return nil, "", err
	}
	doc, err := findDocumentBySubject(dir, subject)
	if err != nil {
		return nil, "", err
	}
	return doc, endpoint, nil
}
