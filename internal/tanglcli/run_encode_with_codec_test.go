//go:build fakecodec

//ff:func feature=cli type=command control=sequence
//ff:what TestRunEncodeCodec — the codec-found path once a "ko" codec is registered: Encode/VerifyTiling errors and the default/explicit output paths
package tanglcli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/spf13/cobra"
)

// TestRunEncodeCodec covers runEncode's branches reachable only once a codec is
// registered for the detected locale: Encode failing, VerifyTiling failing on a
// tiling-broken Doc, the out-flag GetString error, and success with both the
// default output path and an explicit -o path. It registers a single "ko" codec
// (the real en codec owns "en"; doc.RegisterCodec panics on a duplicate locale)
// whose behavior is toggled per subtest via its mutable fields, and only after
// TestRunEncode's own "no codec registered for locale \"ko\"" assertion — this
// file's name sorts after run_encode_test.go, so it always runs later. The fixture
// uses a Korean header so detectLocale returns "ko".
func TestRunEncodeCodec(t *testing.T) {
	codec := &fakeEncodeCodec{}
	doc.RegisterCodec(codec)

	p := filepath.Join(t.TempDir(), "policy.md")
	if err := os.WriteFile(p, []byte("## tangl:주제\n- this document is `x`\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	t.Run("codec Encode error is propagated", func(t *testing.T) {
		codec.encodeErr = errors.New("fakeEncodeCodec: forced Encode failure")
		defer func() { codec.encodeErr = nil }()
		cmd := NewEncodeCmd()
		if err := runEncode(cmd, []string{p}); err == nil || !strings.Contains(err.Error(), "forced Encode failure") {
			t.Fatalf("err = %v, want the forced Encode error", err)
		}
	})

	t.Run("VerifyTiling error is propagated", func(t *testing.T) {
		codec.badTiling = true
		defer func() { codec.badTiling = false }()
		cmd := NewEncodeCmd()
		if err := runEncode(cmd, []string{p}); err == nil || !strings.Contains(err.Error(), "tiling") {
			t.Fatalf("err = %v, want a tiling error", err)
		}
	})

	t.Run("out flag GetString error is propagated", func(t *testing.T) {
		cmd := &cobra.Command{}
		if err := runEncode(cmd, []string{p}); err == nil {
			t.Fatal("expected a flag error for the undefined 'out' flag, got nil")
		}
	})

	t.Run("success with default out path (input path with .tangeul extension)", func(t *testing.T) {
		cmd := NewEncodeCmd()
		if err := runEncode(cmd, []string{p}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		wantOut := strings.TrimSuffix(p, filepath.Ext(p)) + ".tangeul"
		if _, err := os.Stat(wantOut); err != nil {
			t.Fatalf("expected default output at %q: %v", wantOut, err)
		}
	})

	t.Run("success with explicit -o path", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "custom.tangeul")
		cmd := NewEncodeCmd()
		if err := cmd.Flags().Set("out", outPath); err != nil {
			t.Fatalf("Set(out): %v", err)
		}
		if err := runEncode(cmd, []string{p}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := tangeul.Load(outPath); err != nil {
			t.Fatalf("Load(%s): %v", outPath, err)
		}
	})
}
