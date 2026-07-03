//go:build !fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what TestKoE2E — encode a KO surface, check it, and decode it back byte-identically through the real CLI
package tanglcli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestKoE2E drives the full CLI pipeline over a Korean surface exactly as the
// plans/tangl-ko §7 verify commands do, but in-process: registerCodecs wires the
// real KO codec, then `encode test/order-saga-ko.md` produces a .tangeul, `check`
// on it passes with no error, and mode-1 `decode` reproduces the authored source
// byte-for-byte. It is a default-build test — the `fakecodec` build leaves "ko"
// unregistered, so encode there would find no codec.
func TestKoE2E(t *testing.T) {
	registerCodecs()
	src, err := os.ReadFile(filepath.Join("..", "..", "test", "order-saga-ko.md"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	in := filepath.Join(t.TempDir(), "order-saga-ko.md")
	if err := os.WriteFile(in, src, 0o644); err != nil {
		t.Fatalf("write fixture copy: %v", err)
	}
	out := filepath.Join(t.TempDir(), "saga.tangeul")

	enc := NewEncodeCmd()
	if err := enc.Flags().Set("out", out); err != nil {
		t.Fatalf("Set(out): %v", err)
	}
	if err := runEncode(enc, []string{in}); err != nil {
		t.Fatalf("encode: %v", err)
	}

	chk := NewCheckCmd()
	var chkOut, chkErr bytes.Buffer
	chk.SetOut(&chkOut)
	chk.SetErr(&chkErr)
	if err := runCheck(chk, []string{out}); err != nil {
		t.Fatalf("check: %v", err)
	}

	dec := NewDecodeCmd()
	var decOut bytes.Buffer
	dec.SetOut(&decOut)
	if err := runDecode(dec, []string{out}); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !bytes.Equal(decOut.Bytes(), src) {
		t.Errorf("mode-1 decode is not byte-identical to the KO source")
	}
}
