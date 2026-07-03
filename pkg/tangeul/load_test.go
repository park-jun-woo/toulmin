//ff:func feature=tangl type=util control=sequence
//ff:what TestLoad — round-trips a saved Doc, stamps Path, rejects wrong extension and tiling violations
package tangeul

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/stream"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestLoad checks that a saved Doc loads back with the same graph shape and the
// .tangeul path stamped onto AST.Path, that a non-.tangeul path is rejected, and
// that a stream whose Blocks leave a trailing gap (a tiling violation) is refused
// at load time — the .tangeul integrity gate.
func TestLoad(t *testing.T) {
	dir := t.TempDir()

	// Round trip: Save then Load.
	d := buildTransferDoc()
	path := filepath.Join(dir, "transfer.tangeul")
	if err := Save(d, path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.AST.Subject != "transfer" {
		t.Fatalf("subject = %q, want transfer", loaded.AST.Subject)
	}
	if len(loaded.AST.Cases) != 2 {
		t.Fatalf("cases = %d, want 2", len(loaded.AST.Cases))
	}
	if loaded.AST.Path != path {
		t.Fatalf("AST.Path = %q, want %q", loaded.AST.Path, path)
	}

	// Wrong extension is rejected before I/O.
	if _, err := Load(filepath.Join(dir, "transfer.md")); err == nil {
		t.Fatal("Load accepted a .md path")
	}

	// Tiling violation: a Doc whose Blocks cover only part of Source.
	bad := &doc.Doc{
		AST:    &ast.Document{Subject: "x"},
		Source: []byte("----X"),
		Blocks: []doc.Block{{Kind: doc.SectionBlock, Span: doc.Span{Off: 0, Len: 4}, Ref: doc.StmtRef{Sec: codebook.SectionSubject}}},
	}
	b, err := stream.WriteDoc(bad)
	if err != nil {
		t.Fatalf("WriteDoc(bad): %v", err)
	}
	badPath := filepath.Join(dir, "bad.tangeul")
	if err := os.WriteFile(badPath, b, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	_, err = Load(badPath)
	if err == nil {
		t.Fatal("Load accepted a tiling-violating stream")
	}
	if !strings.Contains(err.Error(), "tiling") {
		t.Fatalf("error %q does not mention tiling", err)
	}
}
