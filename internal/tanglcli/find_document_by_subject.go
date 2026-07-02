//ff:func feature=cli type=command control=sequence
//ff:what findDocumentBySubject — recursively scans a directory for the TANGL document whose Subject matches
package tanglcli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// findDocumentBySubject recursively scans dir for ".md" files containing a
// "## tangl:" section marker, parses each, and returns the one Document
// whose Subject equals subject. It errors if none or more than one match.
func findDocumentBySubject(dir, subject string) (*ast.Document, error) {
	var found *ast.Document
	var foundPath string
	walkErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil || !strings.Contains(string(data), "## tangl:") {
			return nil
		}
		doc, err := parser.ParseSource(string(data), path)
		if err != nil || doc.Subject != subject {
			return nil
		}
		if found != nil {
			return fmt.Errorf("tangl: ambiguous subject %q: both %s and %s match", subject, foundPath, path)
		}
		found, foundPath = doc, path
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	if found == nil {
		return nil, fmt.Errorf("tangl: no document with subject %q found under %s", subject, dir)
	}
	return found, nil
}
