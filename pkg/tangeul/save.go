//ff:func feature=tangl type=util control=sequence
//ff:what Save — serialize a Doc to a .tangeul file with an atomic temp+rename write
package tangeul

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/stream"
)

// Save serializes d with stream.WriteDoc and writes the bytes to path. The
// extension is fixed to .tangeul (README contract §5). The write is atomic: bytes
// go to a temporary file in the destination directory and are then renamed over
// path, so a reader never observes a partially written .tangeul file. The temp
// file is removed on any failure before the rename.
func Save(d *doc.Doc, path string) error {
	if filepath.Ext(path) != ".tangeul" {
		return fmt.Errorf("tangeul: Save expects a .tangeul file, got %q", path)
	}
	b, err := stream.WriteDoc(d)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tangeul-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(b); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return err
	}
	return nil
}
