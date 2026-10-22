//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestEncodeDecodeEN — CLI e2e: encode a .md, run check/effects/gen on the .geul, decode back byte-identically
package tanglcli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestEncodeDecodeEN is the Phase 004 L5 end-to-end: it drives the real CLI
// (NewRootCmd, which registers the en codec) over `encode` -> `check` -> `effects`
// -> `gen` -> `decode` for English fixtures and asserts the decoded mode-1 output
// is byte-identical to the original .md. Reaching a codec at all proves the
// registration + locale-detection wiring (README §1). Each subcommand is run on a
// fresh root command; `effects` uses the fixture's own first endpoint.
func TestEncodeDecodeEN(t *testing.T) {
	run := func(t *testing.T, args ...string) string {
		t.Helper()
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return out.String()
	}

	fixtures := []struct{ path, endpoint string }{
		{filepath.Join("..", "..", "pkg", "tangl", "parser", "testdata", "transfer.md"), "transfer"},
		{filepath.Join("..", "..", "..", "test", "order-saga.md"), "place order"},
	}

	for _, f := range fixtures {
		src, err := os.ReadFile(f.path)
		if err != nil {
			t.Fatalf("read %s: %v", f.path, err)
		}
		md := filepath.Join(t.TempDir(), "doc.md")
		if err := os.WriteFile(md, src, 0o644); err != nil {
			t.Fatalf("write %s: %v", md, err)
		}
		tangeulPath := filepath.Join(t.TempDir(), "doc.geul")

		run(t, "encode", md, "-o", tangeulPath)
		if _, err := os.Stat(tangeulPath); err != nil {
			t.Fatalf("%s: encode produced no .geul: %v", f.path, err)
		}
		run(t, "check", tangeulPath)
		run(t, "effects", tangeulPath, f.endpoint)
		run(t, "gen", tangeulPath)

		decoded := run(t, "decode", tangeulPath)
		if !bytes.Equal([]byte(decoded), src) {
			t.Errorf("%s: decode output is not byte-identical to the original .md", f.path)
		}
	}
}
