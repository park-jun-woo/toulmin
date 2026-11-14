//ff:func feature=tangl type=util control=sequence
//ff:what Load — read a .geul file into a Doc (extension check, ReadDoc, tiling gate)
package tangeul

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/stream"
)

// Load reads the .geul file at path and returns its Doc. The extension is
// fixed to .geul (README contract §5); anything else is rejected before any
// I/O. The bytes are parsed by stream.ReadDoc and then gated by doc.VerifyTiling,
// so a tiling violation is a load failure — the .geul file's integrity gate.
// The .geul path is stamped onto d.AST.Path so downstream error messages point
// at the file actually loaded (Path is not stored in the stream; it is metadata).
func Load(path string) (*doc.Doc, error) {
	if filepath.Ext(path) != ".geul" {
		return nil, fmt.Errorf("tangeul: Load expects a .geul file, got %q", path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d, err := stream.ReadDoc(b)
	if err != nil {
		return nil, err
	}
	if err := doc.VerifyTiling(d); err != nil {
		return nil, fmt.Errorf("tangeul: %s: %w", path, err)
	}
	d.AST.Path = path
	return d, nil
}
