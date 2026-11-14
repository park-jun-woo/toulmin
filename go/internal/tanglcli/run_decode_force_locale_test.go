//go:build !fakecodec

//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestRunDecodeForceLocale — CLI e2e: attack forces survive `tangl decode` mode 1 and `--locale en|ko` cross-rendering
package tanglcli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunDecodeForceLocale drives the real CLI over the KO attack-force fixture:
// `encode` → `check` → `decode` (mode 1, byte-identical incl. the authored
// "0.50"/"1.0") → `decode --locale en` (byte-equal to the EN golden) → `encode`
// of that EN rendering → `decode --locale ko` (canonical "0.5의"/"1의" forces).
func TestRunDecodeForceLocale(t *testing.T) {
	run := func(args ...string) string {
		t.Helper()
		cmd := NewRootCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(&out)
		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("%v: %v\n%s", args, err, out.String())
		}
		return out.String()
	}
	fixtures := filepath.Join("..", "..", "pkg", "tangl", "codec", "ko", "testdata")
	koSrc, err := os.ReadFile(filepath.Join(fixtures, "force-ko.md"))
	if err != nil {
		t.Fatal(err)
	}
	enGolden, err := os.ReadFile(filepath.Join(fixtures, "force-en.md"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	koMd, koGeul := filepath.Join(dir, "force-ko.md"), filepath.Join(dir, "force-ko.geul")
	if err := os.WriteFile(koMd, koSrc, 0o644); err != nil {
		t.Fatal(err)
	}
	run("encode", koMd, "-o", koGeul)
	run("check", koGeul)
	if got := run("decode", koGeul); got != string(koSrc) {
		t.Fatalf("mode-1 decode is not byte-identical:\n%s", got)
	}
	enOut := run("decode", "--locale", "en", koGeul)
	if enOut != string(enGolden) {
		t.Fatalf("decode --locale en differs from golden:\n%s", enOut)
	}
	enMd, enGeul := filepath.Join(dir, "force-en.md"), filepath.Join(dir, "force-en.geul")
	if err := os.WriteFile(enMd, []byte(enOut), 0o644); err != nil {
		t.Fatal(err)
	}
	run("encode", enMd, "-o", enGeul)
	koOut := run("decode", "--locale", "ko", enGeul)
	for _, want := range []string{"`신용 충분` 을 0.5의 세기로 공격한다", "`연체 있음` 을 1의 세기로 공격한다", "`신용 부족` 이면 `신용 충분` 무효다", "가중치 0.8 의 반박 규칙이다"} {
		if !strings.Contains(koOut, want) {
			t.Fatalf("decode --locale ko lost %q:\n%s", want, koOut)
		}
	}
}
