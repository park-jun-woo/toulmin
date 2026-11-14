//ff:func feature=tangl type=util control=sequence
//ff:what TestSave — atomic write produces exactly WriteDoc bytes; wrong extension is rejected
package tangeul

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/stream"
)

// TestSave checks that Save writes the .geul file atomically with exactly the
// bytes stream.WriteDoc produces (no partial/extra bytes and no leftover temp
// file in the directory), and that a non-.geul path is rejected before any
// file is created.
func TestSave(t *testing.T) {
	d := buildTransferDoc()
	dir := t.TempDir()
	path := filepath.Join(dir, "transfer.geul")

	if err := Save(d, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	want, err := stream.WriteDoc(d)
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("saved bytes differ from WriteDoc: got %d bytes, want %d bytes", len(got), len(want))
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected only the .geul file, found %d entries (leftover temp?)", len(entries))
	}

	if err := Save(d, filepath.Join(dir, "transfer.md")); err == nil {
		t.Fatal("Save accepted a .md path")
	}
}
