//ff:func feature=cli type=command control=sequence
//ff:what resolveEffectsTarget — resolves the effects command's file+endpoint or subject.endpoint+dir form
package tanglcli

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// resolveEffectsTarget resolves args into a loaded Document and an
// endpoint name. Two args ("<file.tangeul> <endpoint>") load the file
// directly; one arg ("<subject>.<endpoint>") scans dir for a document
// whose Subject matches.
func resolveEffectsTarget(args []string, dir string) (*ast.Document, string, error) {
	if len(args) == 2 {
		d, err := loadExecutable(args[0])
		if err != nil {
			return nil, "", err
		}
		return d.AST, trimBackticks(args[1]), nil
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
