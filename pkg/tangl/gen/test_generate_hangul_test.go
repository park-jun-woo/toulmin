//ff:func feature=tangl type=codegen control=sequence
//ff:what TestGenerateHangulDoc — a Hangul-named document generates distinct, valid, gofmt-clean Go that go build's
package gen

import (
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/parser"
)

// TestGenerateHangulDoc proves the RR romanization fix end to end: a
// document whose case, rule, definition, field and provides names are all
// Hangul must generate a Go file that is valid, gofmt-clean, free of
// identifier collisions, and actually go build's — where before the fix
// every Hangul name collapsed to "_" and the output could not compile. The
// ASCII "logIt" do-target keeps the raw-name Ref path (matched by a stub),
// so this exercises only the ident-generating chokepoints.
func TestGenerateHangulDoc(t *testing.T) {
	src := "" +
		"## tangl:Subject\n- this document is `demo`\n\n" +
		"## tangl:Definitions\n- `기준값` means 5\n\n" +
		"## tangl:Rules\n- `자격` when `점수` is at least 5\n\n" +
		"## tangl:Cases\n" +
		"- in case of `출금 가능`\n" +
		"  - `금액` is required\n" +
		"  - `자격` is a general rule\n" +
		"  - do `logIt` when `자격`\n\n" +
		"- in case of `입금 확인`\n" +
		"  - `승인` is a general rule checking `출금 가능`\n\n" +
		"## tangl:Provides\n" +
		"- provides `출금`\n" +
		"  - `금액` is required\n" +
		"  - run `출금 가능`\n" +
		"- provides `조회`\n" +
		"  - check `입금 확인`\n"
	doc, err := parser.ParseSource(src, "hangul.md")
	if err != nil {
		t.Fatalf("ParseSource: %v", err)
	}
	out, err := Generate(doc)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("generated source is not gofmt-parseable: %v\n--- source ---\n%s", err, out)
	}
	// The two Hangul provides names must romanize to distinct exported idents.
	if !strings.Contains(out, "func Chulgeum(") || !strings.Contains(out, "func Johoe(") {
		t.Fatalf("expected romanized exported endpoints Chulgeum and Johoe in output:\n%s", out)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "demo.go"), []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
	stub := "package demo\n\nimport \"github.com/park-jun-woo/toulmin/pkg/toulmin\"\n\nfunc logIt(ctx toulmin.Context) error { return nil }\n"
	if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte(stub), 0o644); err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	gomod := "module demo\n\ngo 1.18\n\nrequire github.com/park-jun-woo/toulmin v0.0.0\n\nreplace github.com/park-jun-woo/toulmin => " + repoRoot + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off")
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\n%s\n--- generated source ---\n%s", err, cmdOut, out)
	}
}
