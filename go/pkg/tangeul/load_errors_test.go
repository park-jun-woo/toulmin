//ff:func feature=tangl type=util control=sequence
//ff:what TestLoadErrors — Load surfaces os.ReadFile and stream.ReadDoc failures unchanged
package tangeul

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/stream"
)

// TestLoadErrors exercises the two I/O-adjacent error branches that
// TestLoad's round-trip and tiling-violation cases do not reach: a
// .geul path that os.ReadFile cannot open (missing file), and a
// .geul file whose bytes stream.ReadDoc itself rejects (a truncated
// stream, before VerifyTiling is ever consulted).
func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()

	// os.ReadFile fails: the path has the right extension but no file
	// exists there, so Load must propagate the ReadFile error verbatim.
	missing := filepath.Join(dir, "missing.geul")
	if _, err := Load(missing); err == nil {
		t.Fatal("Load accepted a path with no backing file")
	}

	// stream.ReadDoc fails: write a well-formed stream and then drop its
	// last byte, producing an odd-length payload that core.ReadWords (and
	// therefore stream.ReadDoc) rejects before VerifyTiling ever runs.
	b, err := stream.WriteDoc(buildTransferDoc())
	if err != nil {
		t.Fatalf("WriteDoc: %v", err)
	}
	if len(b) == 0 {
		t.Fatalf("WriteDoc produced no bytes")
	}
	corruptPath := filepath.Join(dir, "corrupt.geul")
	if err := os.WriteFile(corruptPath, b[:len(b)-1], 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if _, err := Load(corruptPath); err == nil {
		t.Fatal("Load accepted a stream.ReadDoc-rejected byte stream")
	}
}
