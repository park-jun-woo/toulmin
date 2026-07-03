//ff:func feature=tangl type=util control=sequence
//ff:what TestSaveErrors — Save surfaces WriteDoc, temp-file, write and rename failures unchanged
package tangeul

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

// TestSaveErrors exercises the atomic-write error branches that TestSave's
// happy-path and extension-rejection cases do not reach: stream.WriteDoc
// failing on a Doc with no AST, os.CreateTemp failing because the
// destination directory does not exist, the temp file's Write failing
// (simulated via RLIMIT_FSIZE so it deterministically returns EFBIG instead
// of relying on real disk exhaustion), and os.Rename failing because the
// destination path is already occupied by a directory. Every failure must
// leave no dangling .tangeul file at path.
func TestSaveErrors(t *testing.T) {
	dir := t.TempDir()

	// stream.WriteDoc fails: WriteDoc's own nil-Doc/nil-AST guard, reached
	// here through Save with a nil *doc.Doc.
	if err := Save(nil, filepath.Join(dir, "nil-ast.tangeul")); err == nil {
		t.Fatal("Save accepted a nil Doc")
	}

	// os.CreateTemp fails: the destination directory does not exist, so no
	// temp file can be created in it.
	missingDir := filepath.Join(dir, "no-such-subdir", "x.tangeul")
	if err := Save(buildTransferDoc(), missingDir); err == nil {
		t.Fatal("Save accepted a path in a non-existent directory")
	}

	// os.Rename fails: the destination is already an existing directory, so
	// renaming the finished temp file onto it is rejected by the OS. The
	// temp file must be cleaned up rather than left behind — a dedicated
	// subdirectory isolates the check to a single, non-iterating count.
	renameDir := filepath.Join(dir, "rename-case")
	if err := os.Mkdir(renameDir, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	occupied := filepath.Join(renameDir, "occupied.tangeul")
	if err := os.Mkdir(occupied, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	if err := Save(buildTransferDoc(), occupied); err == nil {
		t.Fatal("Save accepted a path already occupied by a directory")
	}
	entries, err := os.ReadDir(renameDir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("rename-case dir has %d entries, want 1 (leftover temp file after Rename failure)", len(entries))
	}

	// tmp.Write fails: lower RLIMIT_FSIZE below the size of the encoded
	// stream so the temp file's Write syscall deterministically returns
	// EFBIG. SIGXFSZ is ignored first so the process observes the error
	// instead of being killed by the default disposition, and both the
	// signal disposition and the limit are restored before the test ends.
	if runtime.GOOS != "linux" {
		t.Skip("RLIMIT_FSIZE write-failure simulation is Linux-specific")
	}
	signal.Ignore(syscall.SIGXFSZ)
	defer signal.Reset(syscall.SIGXFSZ)
	var orig syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_FSIZE, &orig); err != nil {
		t.Fatalf("Getrlimit: %v", err)
	}
	defer syscall.Setrlimit(syscall.RLIMIT_FSIZE, &orig)
	if err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &syscall.Rlimit{Cur: 1, Max: orig.Max}); err != nil {
		t.Fatalf("Setrlimit: %v", err)
	}
	writePath := filepath.Join(dir, "write-fail.tangeul")
	err = Save(buildTransferDoc(), writePath)
	if rerr := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &orig); rerr != nil {
		t.Fatalf("Setrlimit restore: %v", rerr)
	}
	if err == nil {
		t.Fatal("Save succeeded despite RLIMIT_FSIZE forcing the temp write to fail")
	}
	if _, statErr := os.Stat(writePath); !os.IsNotExist(statErr) {
		t.Fatalf("Save left a file at %s after a write failure", writePath)
	}
}
