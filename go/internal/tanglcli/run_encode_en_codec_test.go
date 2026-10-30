//ff:func feature=cli type=command control=sequence
//ff:what TestRunEncodeENCodec — with the real EN codec: parse errors, the missing out flag, and default/explicit output paths
package tanglcli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul"
	"github.com/spf13/cobra"
)

// TestRunEncodeENCodec drives runEncode's post-lookup branches in the default
// build through the real EN codec (registered by registerCodecs in every build):
// an EN surface the parser rejects propagates the Encode error; a command without
// an "out" flag propagates GetString's error; with no -o the .geul lands next to
// the input, and an explicit -o path is honored and loads back to the same subject.
func TestRunEncodeENCodec(t *testing.T) {
	registerCodecs()
	write := func(t *testing.T, body string) string {
		t.Helper()
		p := filepath.Join(t.TempDir(), "policy.md")
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		return p
	}
	const valid = "## tangl:Subject\n- this document is `x`\n\n" +
		"## tangl:Rules\n- `qualifies` when `score` is at least 5\n\n" +
		"## tangl:Cases\n- in case of `gate`\n  - `qualifies` is a general rule\n"

	t.Run("EN parse error is propagated", func(t *testing.T) {
		p := write(t, "## tangl:Subject\n- this is not a subject sentence\n")
		if err := runEncode(NewEncodeCmd(), []string{p}); err == nil {
			t.Fatal("expected the EN codec to reject a malformed subject")
		}
		if _, err := os.Stat(strings.TrimSuffix(p, ".md") + ".geul"); !os.IsNotExist(err) {
			t.Fatalf("no .geul may be written on a parse error: %v", err)
		}
	})

	t.Run("missing out flag is propagated", func(t *testing.T) {
		p := write(t, valid)
		err := runEncode(&cobra.Command{}, []string{p})
		if err == nil || !strings.Contains(err.Error(), "out") {
			t.Fatalf("err = %v, want the undefined 'out' flag error", err)
		}
	})

	t.Run("default output path replaces .md with .geul", func(t *testing.T) {
		p := write(t, valid)
		if err := runEncode(NewEncodeCmd(), []string{p}); err != nil {
			t.Fatalf("runEncode: %v", err)
		}
		d, err := tangeul.Load(strings.TrimSuffix(p, ".md") + ".geul")
		if err != nil {
			t.Fatalf("Load default output: %v", err)
		}
		if d.AST.Subject != "x" || d.Locale != "en" {
			t.Fatalf("subject=%q locale=%q, want x/en", d.AST.Subject, d.Locale)
		}
	})

	t.Run("explicit -o path is honored", func(t *testing.T) {
		p := write(t, valid)
		out := filepath.Join(t.TempDir(), "custom.geul")
		cmd := NewEncodeCmd()
		if err := cmd.Flags().Set("out", out); err != nil {
			t.Fatalf("Set(out): %v", err)
		}
		if err := runEncode(cmd, []string{p}); err != nil {
			t.Fatalf("runEncode: %v", err)
		}
		if _, err := os.Stat(strings.TrimSuffix(p, ".md") + ".geul"); !os.IsNotExist(err) {
			t.Fatalf("default path must stay unused when -o is set: %v", err)
		}
		if d, err := tangeul.Load(out); err != nil || d.AST.Subject != "x" {
			t.Fatalf("Load(-o output): %v", err)
		}
	})
}
