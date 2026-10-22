//ff:func feature=cli type=command control=sequence
//ff:what findDocumentBySubject — recursively scans a directory for the .geul document whose Subject matches
package tanglcli

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// findDocumentBySubject recursively scans dir for ".geul" files, loads each,
// and returns the one Document whose Subject equals subject. It errors if none or
// more than one match. Files that fail to load are skipped (not the target).
func findDocumentBySubject(dir, subject string) (*ast.Document, error) {
	var found *ast.Document
	var foundPath string
	walkErr := filepath.WalkDir(dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil || de.IsDir() || filepath.Ext(path) != ".geul" {
			return err
		}
		d, err := tangeul.Load(path)
		if err != nil || d.AST.Subject != subject {
			return nil
		}
		if found != nil {
			return fmt.Errorf("tangl: ambiguous subject %q: both %s and %s match", subject, foundPath, path)
		}
		found, foundPath = d.AST, path
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
