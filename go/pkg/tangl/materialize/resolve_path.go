//ff:func feature=tangl type=engine control=sequence
//ff:what resolvePath — resolve a definition file path (absolute as-is, relative to the document's directory)
package materialize

import (
	"fmt"
	"path/filepath"
)

// resolvePath cleans an absolute sourcePath, or joins a relative one onto the
// directory of documentPath; a relative path without a document path is an error.
func resolvePath(documentPath, sourcePath string) (string, error) {
	if filepath.IsAbs(sourcePath) {
		return filepath.Clean(sourcePath), nil
	}
	if documentPath == "" {
		return "", fmt.Errorf("relative file path requires Document.Path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(documentPath), sourcePath)), nil
}
